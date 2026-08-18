package app

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ohnishi/feed2mail/internal/feed"
	"github.com/ohnishi/feed2mail/internal/store"
)

// runner は同じ状態ファイルを使い回して複数回の実行を再現します。
type runner struct {
	t         *testing.T
	statePath string
	seenPath  string
	notifier  *fakeNotifier
	now       time.Time
	retention time.Duration
}

func newRunner(t *testing.T, subs []store.Subscription) *runner {
	t.Helper()
	dir := t.TempDir()
	r := &runner{
		t:         t,
		statePath: filepath.Join(dir, "fetched.jsonl"),
		seenPath:  filepath.Join(dir, "seen.jsonl"),
		notifier:  &fakeNotifier{},
		now:       fixedNow,
		retention: 30 * 24 * time.Hour,
	}
	if err := store.SaveSubscriptions(r.statePath, subs); err != nil {
		t.Fatalf("SaveSubscriptions() error = %v", err)
	}
	return r
}

// run は 1 回分の実行を行い、その回に送信された本文を返します。未送信なら空文字列。
func (r *runner) run(feeds map[string]*feed.Feed) string {
	r.t.Helper()
	before := len(r.notifier.calls)

	a := &App{
		Fetcher:         &fakeFetcher{feeds: feeds},
		Notifier:        r.notifier,
		StatePath:       r.statePath,
		SeenPath:        r.seenPath,
		SeenRetention:   r.retention,
		Subject:         "feed更新通知",
		ExcludePrefixes: []string{"https://anond.hatelabo.jp"},
		Sleep:           func(time.Duration) {},
		Now:             func() time.Time { return r.now },
	}
	if err := a.Run(context.Background()); err != nil {
		r.t.Fatalf("Run() error = %v", err)
	}

	if len(r.notifier.calls) == before {
		return ""
	}
	return r.notifier.calls[len(r.notifier.calls)-1]
}

func (r *runner) seen() *store.SeenStore {
	r.t.Helper()
	s, err := store.LoadSeen(r.seenPath)
	if err != nil {
		r.t.Fatalf("LoadSeen() error = %v", err)
	}
	return s
}

const feedURL = "https://a.example.com/feed"

func oneFeed(items ...feed.Item) map[string]*feed.Feed {
	return map[string]*feed.Feed{feedURL: {Items: items}}
}

func subs() []store.Subscription {
	return []store.Subscription{{Title: "A", URL: feedURL, Fetched: hoursAgo(24)}}
}

// 移行直後（既読記録なし）は Fetched より新しいものだけを通知し、残りは既読として取り込む。
func TestFirstRunSeedsWithoutFloodingNotification(t *testing.T) {
	r := newRunner(t, subs())

	body := r.run(oneFeed(
		feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		feed.Item{Title: "古い1", Link: "https://a.example.com/old1", Published: hoursAgo(100)},
		feed.Item{Title: "古い2", Link: "https://a.example.com/old2", Published: hoursAgo(200)},
	))

	if !strings.Contains(body, "/new") {
		t.Errorf("body missing new item: %q", body)
	}
	if strings.Contains(body, "/old1") || strings.Contains(body, "/old2") {
		t.Errorf("seeding leaked old items into notification: %q", body)
	}
	// 通知しなかった古い記事も既読として取り込み、次回に浮かび上がらせない。
	if n := r.seen().Len(); n != 3 {
		t.Errorf("seen count = %d, want 3", n)
	}
}

// 既読になったリンクは、次回以降に再通知しない。
func TestSeenLinkIsNotNotifiedTwice(t *testing.T) {
	r := newRunner(t, subs())
	items := oneFeed(feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)})

	if body := r.run(items); !strings.Contains(body, "/new") {
		t.Fatalf("first run body = %q, want the new item", body)
	}
	if body := r.run(items); body != "" {
		t.Errorf("second run notified again: %q", body)
	}
}

// 既読方式の主眼。Fetched より古い公開時刻でも、後から現れた記事なら通知する。
func TestBackdatedItemIsNotifiedOnceSeeded(t *testing.T) {
	r := newRunner(t, subs())

	// 1 回目でこのフィードの既読記録ができ、Fetched は hoursAgo(1) まで進む。
	r.run(oneFeed(feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)}))

	// 2 回目に、Fetched より古い公開時刻の記事が新たに現れる。
	body := r.run(oneFeed(
		feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		feed.Item{Title: "遡って追加", Link: "https://a.example.com/backdated", Published: hoursAgo(5)},
	))

	if !strings.Contains(body, "/backdated") {
		t.Errorf("backdated item was missed: %q", body)
	}
}

// 保持期間より古い記事は、既読記録が無くても新着として扱わない。
func TestItemsOlderThanRetentionAreNotNotified(t *testing.T) {
	r := newRunner(t, subs())
	r.run(oneFeed(feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)}))

	body := r.run(oneFeed(
		feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		feed.Item{Title: "大昔", Link: "https://a.example.com/ancient", Published: r.now.Add(-31 * 24 * time.Hour)},
	))

	if strings.Contains(body, "/ancient") {
		t.Errorf("item older than retention was notified: %q", body)
	}
}

// 保持期間を過ぎた既読記録は捨てる。
func TestSeenLinksArePrunedAfterRetention(t *testing.T) {
	r := newRunner(t, subs())
	r.run(oneFeed(feed.Item{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)}))
	if n := r.seen().Len(); n != 1 {
		t.Fatalf("seen count = %d, want 1", n)
	}

	// 保持期間を越えて時計を進めると、記録は落ちる。
	r.now = r.now.Add(31 * 24 * time.Hour)
	r.run(oneFeed())

	if n := r.seen().Len(); n != 0 {
		t.Errorf("seen count = %d, want 0 after pruning", n)
	}
}

// 除外対象のリンクは通知もせず、既読にも記録しない（状態ファイルを太らせないため）。
func TestExcludedLinksAreNotRecorded(t *testing.T) {
	r := newRunner(t, subs())

	body := r.run(oneFeed(
		feed.Item{Title: "匿名", Link: "https://anond.hatelabo.jp/20260818", Published: hoursAgo(1)},
		feed.Item{Title: "通常", Link: "https://a.example.com/new", Published: hoursAgo(1)},
	))

	if strings.Contains(body, "anond.hatelabo.jp") {
		t.Errorf("excluded link was notified: %q", body)
	}
	if r.seen().Has("https://anond.hatelabo.jp/20260818") {
		t.Error("excluded link was recorded as seen")
	}
}

// 新着が無い実行でも既読記録は保存する。保存しないと次回に再通知してしまう。
func TestSeenStateIsSavedEvenWithoutNotification(t *testing.T) {
	r := newRunner(t, subs())
	// Fetched より古いので通知されないが、既読としては取り込まれる。
	if body := r.run(oneFeed(feed.Item{Title: "古い", Link: "https://a.example.com/old", Published: hoursAgo(100)})); body != "" {
		t.Fatalf("unexpected notification: %q", body)
	}
	if !r.seen().Has("https://a.example.com/old") {
		t.Error("seen state was not saved when nothing was notified")
	}
}
