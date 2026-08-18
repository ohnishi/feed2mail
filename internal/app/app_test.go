package app

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ohnishi/feed2mail/internal/feed"
	"github.com/ohnishi/feed2mail/internal/store"
)

var fixedNow = time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)

// fakeFetcher は URL ごとに固定の応答を返します。
type fakeFetcher struct {
	feeds  map[string]*feed.Feed
	errs   map[string]error
	calls  []string
	sleeps []time.Duration
}

func (f *fakeFetcher) Fetch(_ context.Context, url string) (*feed.Feed, error) {
	f.calls = append(f.calls, url)
	if err, ok := f.errs[url]; ok {
		return nil, err
	}
	if fd, ok := f.feeds[url]; ok {
		return fd, nil
	}
	return &feed.Feed{}, nil
}

// fakeNotifier は送信内容を記録するだけの Notifier です。
type fakeNotifier struct {
	calls []string
	err   error
}

func (n *fakeNotifier) Notify(_ context.Context, subject, body string) error {
	if n.err != nil {
		return n.err
	}
	n.calls = append(n.calls, subject+"\n"+body)
	return nil
}

func hoursAgo(h int) time.Time { return fixedNow.Add(-time.Duration(h) * time.Hour) }

// newTestApp は状態ファイルを書き出したうえで App を組み立てます。
func newTestApp(t *testing.T, subs []store.Subscription, fetcher *fakeFetcher, notifier *fakeNotifier) (*App, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fetched.jsonl")
	if err := store.SaveSubscriptions(path, subs); err != nil {
		t.Fatalf("SaveSubscriptions() error = %v", err)
	}
	return &App{
		Fetcher:         fetcher,
		Notifier:        notifier,
		StatePath:       path,
		SeenPath:        filepath.Join(filepath.Dir(path), "seen.jsonl"),
		SeenRetention:   30 * 24 * time.Hour,
		Subject:         "feed更新通知",
		ExcludePrefixes: []string{"https://anond.hatelabo.jp"},
		PoliteDelay:     2 * time.Second,
		Sleep:           func(d time.Duration) { fetcher.sleeps = append(fetcher.sleeps, d) },
		Now:             func() time.Time { return fixedNow },
	}, path
}

func TestRunNotifiesNewItemsAndAdvancesFetched(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
			{Title: "既読", Link: "https://a.example.com/old", Published: hoursAgo(48)},
		}},
	}}
	notifier := &fakeNotifier{}
	a, path := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("Notify called %d times, want 1", len(notifier.calls))
	}
	body := notifier.calls[0]
	if !strings.Contains(body, "https://a.example.com/new") {
		t.Errorf("body missing new item: %q", body)
	}
	if strings.Contains(body, "https://a.example.com/old") {
		t.Errorf("body contains already-fetched item: %q", body)
	}

	saved, err := store.LoadSubscriptions(path)
	if err != nil {
		t.Fatalf("LoadSubscriptions() error = %v", err)
	}
	if !saved[0].Fetched.Equal(hoursAgo(1)) {
		t.Errorf("Fetched = %v, want %v", saved[0].Fetched, hoursAgo(1))
	}
}

func TestRunSkipsNotificationWhenNothingNew(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "既読", Link: "https://a.example.com/old", Published: hoursAgo(48)},
		}},
	}}
	notifier := &fakeNotifier{}
	a, _ := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(notifier.calls) != 0 {
		t.Errorf("Notify called %d times, want 0", len(notifier.calls))
	}
}

// 未来日付のアイテムで Fetched が飛ぶと以降の新着を恒久的に取りこぼす。
func TestRunIgnoresFutureDatedItems(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "未来", Link: "https://a.example.com/future", Published: fixedNow.Add(240 * time.Hour)},
			{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		}},
	}}
	notifier := &fakeNotifier{}
	a, path := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	saved, _ := store.LoadSubscriptions(path)
	if !saved[0].Fetched.Equal(hoursAgo(1)) {
		t.Errorf("Fetched = %v, want %v (未来日付に引きずられている)", saved[0].Fetched, hoursAgo(1))
	}
	if strings.Contains(notifier.calls[0], "future") {
		t.Errorf("future-dated item was notified: %q", notifier.calls[0])
	}
}

func TestRunExcludesConfiguredPrefixes(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "匿名", Link: "https://anond.hatelabo.jp/20260818", Published: hoursAgo(1)},
			{Title: "通常", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		}},
	}}
	notifier := &fakeNotifier{}
	a, _ := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if strings.Contains(notifier.calls[0], "anond.hatelabo.jp") {
		t.Errorf("excluded link was notified: %q", notifier.calls[0])
	}
}

func TestRunDeduplicatesLinksAcrossFeeds(t *testing.T) {
	shared := feed.Item{Title: "共有", Link: "https://shared.example.com/x", Published: hoursAgo(1)}
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{shared}},
		"https://b.example.com/feed": {Items: []feed.Item{shared}},
	}}
	notifier := &fakeNotifier{}
	a, _ := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
		{Title: "B", URL: "https://b.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if n := strings.Count(notifier.calls[0], "https://shared.example.com/x"); n != 1 {
		t.Errorf("shared link appeared %d times, want 1", n)
	}
}

// 1 つのフィードの失敗で全体を止めず、失敗したフィードの Fetched は据え置く。
func TestRunContinuesOnFetchFailure(t *testing.T) {
	fetcher := &fakeFetcher{
		feeds: map[string]*feed.Feed{
			"https://b.example.com/feed": {Items: []feed.Item{
				{Title: "新着", Link: "https://b.example.com/new", Published: hoursAgo(1)},
			}},
		},
		errs: map[string]error{"https://a.example.com/feed": errors.New("boom")},
	}
	notifier := &fakeNotifier{}
	a, path := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
		{Title: "B", URL: "https://b.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	body := notifier.calls[0]
	if !strings.Contains(body, "取得に失敗したフィード:") || !strings.Contains(body, "  - A") {
		t.Errorf("body missing failure section: %q", body)
	}
	if !strings.Contains(body, "https://b.example.com/new") {
		t.Errorf("body missing item from healthy feed: %q", body)
	}

	saved, _ := store.LoadSubscriptions(path)
	if !saved[0].Fetched.Equal(hoursAgo(24)) {
		t.Errorf("failed feed Fetched = %v, want unchanged %v", saved[0].Fetched, hoursAgo(24))
	}
}

// 通知に失敗したら状態を進めず、次回に取りこぼしを拾い直せるようにする。
func TestRunDoesNotAdvanceStateWhenNotifyFails(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "新着", Link: "https://a.example.com/new", Published: hoursAgo(1)},
		}},
	}}
	notifier := &fakeNotifier{err: errors.New("smtp down")}
	a, path := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err == nil {
		t.Fatal("Run() error = nil, want error")
	}

	saved, _ := store.LoadSubscriptions(path)
	if !saved[0].Fetched.Equal(hoursAgo(24)) {
		t.Errorf("Fetched = %v, want unchanged %v", saved[0].Fetched, hoursAgo(24))
	}
}

// 公開時刻が判定できないアイテムは通知対象にしない。
func TestRunSkipsItemsWithoutPublishedTime(t *testing.T) {
	fetcher := &fakeFetcher{feeds: map[string]*feed.Feed{
		"https://a.example.com/feed": {Items: []feed.Item{
			{Title: "日付なし", Link: "https://a.example.com/undated"},
		}},
	}}
	notifier := &fakeNotifier{}
	a, _ := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed", Fetched: hoursAgo(24)},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(notifier.calls) != 0 {
		t.Errorf("Notify called %d times, want 0", len(notifier.calls))
	}
}

// フィード間には PoliteDelay を挟むが、先頭の 1 件では待たない。
func TestRunSleepsBetweenFeedsOnly(t *testing.T) {
	fetcher := &fakeFetcher{}
	notifier := &fakeNotifier{}
	a, _ := newTestApp(t, []store.Subscription{
		{Title: "A", URL: "https://a.example.com/feed"},
		{Title: "B", URL: "https://b.example.com/feed"},
		{Title: "C", URL: "https://c.example.com/feed"},
	}, fetcher, notifier)

	if err := a.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(fetcher.sleeps) != 2 {
		t.Errorf("sleep called %d times, want 2", len(fetcher.sleeps))
	}
	for _, d := range fetcher.sleeps {
		if d != 2*time.Second {
			t.Errorf("sleep duration = %v, want 2s", d)
		}
	}
}
