// Command feed は購読中の RSS/Atom フィードを巡回し、新着記事をメールで通知します。
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ohnishi/feed/internal/app"
	"github.com/ohnishi/feed/internal/feed"
	"github.com/ohnishi/feed/internal/notify"
	"github.com/ohnishi/feed/internal/store"
)

var subscriptions = []store.Subscription{
	{Title: "はてブ（総合）", URL: "https://b.hatena.ne.jp/hotentry.rss"},
	{Title: "はてブ（テクノロジー）", URL: "https://b.hatena.ne.jp/hotentry/it.rss"},

	{Title: "Zenn", URL: "https://zenn.dev/feed"},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml"},

	{Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf"},
	{Title: "AI Watch", URL: "https://ai.watch.impress.co.jp/data/rss/1.0/aiw/feed.rdf"},

	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf"},
	{Title: "くるねこ大和", URL: "https://rssblog.ameba.jp/kuru0214neko/rss20.xml"},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490"},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745"},

	{Title: "物販NAVI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMO-lX7Y0K7c-BaV9RatboQ"},
	{Title: "中古カメラ現物投資家 船田ひろし", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcph2C76gdSHXUbd3wooNfQ"},
	{Title: "脱サラした男のネットショップな日常", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfCkw1ugUhc9Nw_wu1cKHpw"},
	{Title: "かわしま＠中国輸入", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCW_Ey8M8I3o1JCB2OyCYzHA"},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA"},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA"},
}

// excludePrefixes はこの接頭辞で始まるリンクを通知対象から外します。
var excludePrefixes = []string{
	"https://anond.hatelabo.jp",
	"https://www.youtube.com/shorts/",
}

const (
	maxRetry       = 5
	statePath      = "cache/fetched.jsonl"
	mailSubject    = "feed更新通知"
	mailFrom       = "feed@resend.dev"
	mailTo         = "notify@example.com"
	politeDelay    = 2 * time.Second
	newSubLookback = 24 * time.Hour
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("feed: %v", err)
	}
}

func run() error {
	if err := syncSubscriptions(statePath); err != nil {
		return err
	}

	application := &app.App{
		Fetcher:         feed.NewGofeedFetcher(maxRetry),
		Notifier:        notify.NewResend(os.Getenv("MAIL_APIKEY"), mailFrom, []string{mailTo}, maxRetry),
		StatePath:       statePath,
		Subject:         mailSubject,
		ExcludePrefixes: excludePrefixes,
		PoliteDelay:     politeDelay,
	}

	return application.Run(context.Background())
}

// syncSubscriptions は購読リストの定義を状態ファイルへ反映します。
// 既知のフィードは取得済み時刻を引き継ぎ、Title は購読リスト側の定義を優先します。
// 定義から外れたフィードは状態ファイルからも落ちます。
func syncSubscriptions(path string) error {
	saved, err := store.LoadSubscriptions(path)
	if err != nil {
		return err
	}

	fetchedByURL := make(map[string]time.Time, len(saved))
	for _, s := range saved {
		fetchedByURL[s.URL] = s.Fetched
	}

	synced := make([]store.Subscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		if fetched, ok := fetchedByURL[s.URL]; ok {
			s.Fetched = fetched
		} else {
			// 新規購読は直近 1 日分だけを初回通知の対象にする。
			s.Fetched = time.Now().Add(-newSubLookback)
		}
		synced = append(synced, s)
	}

	return store.SaveSubscriptions(path, synced)
}
