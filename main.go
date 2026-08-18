// Command feed は購読中の RSS/Atom フィードを巡回し、新着記事をメールで通知します。
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ohnishi/feed/internal/app"
	"github.com/ohnishi/feed/internal/config"
	"github.com/ohnishi/feed/internal/feed"
	"github.com/ohnishi/feed/internal/notify"
	"github.com/ohnishi/feed/internal/store"
)

func main() {
	configPath := flag.String("config", config.DefaultPath, "path to the config file")
	flag.Parse()

	if err := run(*configPath); err != nil {
		log.Fatalf("feed: %v", err)
	}
}

func run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if err := syncSubscriptions(cfg); err != nil {
		return err
	}

	application := &app.App{
		Fetcher:         feed.NewGofeedFetcher(cfg.MaxRetry),
		Notifier:        notify.NewResend(cfg.APIKey, cfg.Mail.From, cfg.Mail.To, cfg.MaxRetry),
		StatePath:       cfg.StatePath,
		SeenPath:        cfg.SeenPath,
		SeenRetention:   cfg.SeenRetention,
		Subject:         cfg.Mail.Subject,
		ExcludePrefixes: cfg.ExcludePrefixes,
		PoliteDelay:     cfg.PoliteDelay,
	}

	return application.Run(context.Background())
}

// syncSubscriptions は設定ファイルの購読リストを状態ファイルへ反映します。
// 既知のフィードは取得済み時刻を引き継ぎ、Title は設定ファイル側の定義を優先します。
// 設定から外れたフィードは状態ファイルからも落ちます。
func syncSubscriptions(cfg config.Config) error {
	saved, err := store.LoadSubscriptions(cfg.StatePath)
	if err != nil {
		return err
	}

	fetchedByURL := make(map[string]time.Time, len(saved))
	for _, s := range saved {
		fetchedByURL[s.URL] = s.Fetched
	}

	synced := make([]store.Subscription, 0, len(cfg.Subscriptions))
	for _, s := range cfg.Subscriptions {
		sub := store.Subscription{Title: s.Title, URL: s.URL}
		if fetched, ok := fetchedByURL[s.URL]; ok {
			sub.Fetched = fetched
		} else {
			sub.Fetched = time.Now().Add(-cfg.NewSubLookback)
		}
		synced = append(synced, sub)
	}

	return store.SaveSubscriptions(cfg.StatePath, synced)
}
