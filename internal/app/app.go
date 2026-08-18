// Package app はフィード取得から通知・保存までの一連の流れを組み立てます。
// 入出力はすべてインタフェース経由なので、テストでは実ネットワークもファイルも使いません。
package app

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ohnishi/feed/internal/feed"
	"github.com/ohnishi/feed/internal/notify"
	"github.com/ohnishi/feed/internal/render"
	"github.com/ohnishi/feed/internal/store"
)

// App は 1 回の実行に必要な依存をまとめたものです。
type App struct {
	// Fetcher はフィードを取得します。
	Fetcher feed.Fetcher
	// Notifier は本文を送信します。
	Notifier notify.Notifier
	// StatePath は購読状態を保存する JSONL のパスです。
	StatePath string
	// Subject は通知メールの件名です。
	Subject string
	// ExcludePrefixes はこの接頭辞で始まるリンクを通知対象から除外します。
	ExcludePrefixes []string
	// PoliteDelay はフィード間の待機時間です。
	// YouTube のフィードは短時間に連続アクセスすると 404/500 を返すことがあるため、間隔を空ける。
	PoliteDelay time.Duration

	// Sleep と Now はテストから差し替えるための注入点です。ゼロ値なら実時間を使います。
	Sleep func(time.Duration)
	Now   func() time.Time
}

func (a *App) sleep(d time.Duration) {
	if a.Sleep != nil {
		a.Sleep(d)
		return
	}
	time.Sleep(d)
}

func (a *App) now() time.Time {
	if a.Now != nil {
		return a.Now()
	}
	return time.Now()
}

// Run は購読状態を読み込み、各フィードの新着を収集して通知し、状態を書き戻します。
func (a *App) Run(ctx context.Context) error {
	subscriptions, err := store.LoadSubscriptions(a.StatePath)
	if err != nil {
		return err
	}

	sections, failed := a.collect(ctx, subscriptions)

	body := render.Body(sections, failed)
	if body == "" {
		// 新着も失敗もない場合は空メールを送らない。
		log.Print("no new items; skip notification")
		return store.SaveSubscriptions(a.StatePath, subscriptions)
	}

	if err := a.Notifier.Notify(ctx, a.Subject, body); err != nil {
		// 通知に失敗したときは状態を進めない。次回の実行で取りこぼしを拾い直せる。
		return err
	}

	return store.SaveSubscriptions(a.StatePath, subscriptions)
}

// collect は各フィードを順に取得し、新着セクションと失敗フィード名を返します。
// subscriptions の Fetched は取得できたフィードについてのみ更新されます。
func (a *App) collect(ctx context.Context, subscriptions []store.Subscription) ([]render.Section, []string) {
	var sections []render.Section
	var failed []string
	// 同じリンクが複数のフィードに現れても 1 度しか通知しない。
	seen := make(map[string]struct{})

	for i, subscription := range subscriptions {
		if i > 0 {
			a.sleep(a.PoliteDelay)
		}

		fetched, err := a.Fetcher.Fetch(ctx, subscription.URL)
		if err != nil {
			// 1 つのフィードの失敗で全体を止めず、取得できた分だけ通知する。
			// Fetched も更新しないので、次回の実行で取りこぼしを拾い直せる。
			log.Printf("skip %s: %v", subscription.Title, err)
			failed = append(failed, subscription.Title)
			continue
		}

		section, latest := a.selectItems(subscription, fetched, seen)
		if len(section.Items) > 0 {
			sections = append(sections, section)
		}
		if subscription.Fetched.Before(latest) {
			subscriptions[i].Fetched = latest
		}
	}

	return sections, failed
}

// selectItems は 1 フィードから未通知の新着を抜き出し、併せて最新の公開時刻を返します。
// 通知対象に選んだリンクは seen に記録されます。
func (a *App) selectItems(subscription store.Subscription, fetched *feed.Feed, seen map[string]struct{}) (render.Section, time.Time) {
	now := a.now()
	section := render.Section{Title: subscription.Title}
	var latest time.Time

	for _, item := range fetched.Items {
		if !item.HasPublished() {
			continue
		}
		// 未来日付のアイテムで Fetched が先に進むと、以降の新着を恒久的に
		// 取りこぼすため、現在時刻を超える公開時刻は latest に反映しない。
		if item.Published.After(now) {
			continue
		}

		// 除外対象でも latest には反映する。次回そのアイテムを再評価しないため。
		if latest.Before(item.Published) {
			latest = item.Published
		}

		if !subscription.Fetched.Before(item.Published) {
			continue
		}
		if _, dup := seen[item.Link]; dup {
			continue
		}
		if a.excluded(item.Link) {
			continue
		}

		section.Items = append(section.Items, render.Item{Title: item.Title, Link: item.Link})
		seen[item.Link] = struct{}{}
	}

	return section, latest
}

func (a *App) excluded(link string) bool {
	for _, prefix := range a.ExcludePrefixes {
		if strings.HasPrefix(link, prefix) {
			return true
		}
	}
	return false
}
