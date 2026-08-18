package feed

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"><channel><title>sample</title>
<item><title>記事A</title><link>https://example.com/a</link><pubDate>Mon, 17 Aug 2026 09:00:00 +0900</pubDate></item>
<item><title>日付なし</title><link>https://example.com/b</link></item>
</channel></rss>`

// newTestFetcher は待機を記録するだけの sleep を持つ Fetcher を返します。
func newTestFetcher(retry int, sleeps *[]time.Duration) *GofeedFetcher {
	return NewGofeedFetcher(retry, WithSleep(func(d time.Duration) { *sleeps = append(*sleeps, d) }))
}

func TestFetchParsesItems(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(sampleRSS))
	}))
	defer srv.Close()

	var sleeps []time.Duration
	got, err := newTestFetcher(3, &sleeps).Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(got.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(got.Items))
	}
	if got.Items[0].Title != "記事A" || got.Items[0].Link != "https://example.com/a" {
		t.Errorf("Items[0] = %+v", got.Items[0])
	}
	if !got.Items[0].HasPublished() {
		t.Error("Items[0].HasPublished() = false, want true")
	}
	if got.Items[1].HasPublished() {
		t.Error("Items[1].HasPublished() = true, want false (公開時刻が無い)")
	}
	if len(sleeps) != 0 {
		t.Errorf("sleep called %d times on success, want 0", len(sleeps))
	}
}

func TestFetchRetriesWithIncreasingBackoff(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(sampleRSS))
	}))
	defer srv.Close()

	var sleeps []time.Duration
	if _, err := newTestFetcher(5, &sleeps).Fetch(context.Background(), srv.URL); err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	want := []time.Duration{DefaultRetryInterval, 2 * DefaultRetryInterval}
	if len(sleeps) != len(want) {
		t.Fatalf("sleeps = %v, want %v", sleeps, want)
	}
	for i := range want {
		if sleeps[i] != want[i] {
			t.Errorf("sleeps[%d] = %v, want %v", i, sleeps[i], want[i])
		}
	}
}

func TestFetchGivesUpAfterRetries(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	var sleeps []time.Duration
	if _, err := newTestFetcher(3, &sleeps).Fetch(context.Background(), srv.URL); err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	// 最後の試行後は待たない。
	if len(sleeps) != 2 {
		t.Errorf("sleeps = %v, want 2 entries", sleeps)
	}
}

// retry が 1 未満でも必ず 1 回は試行する。
func TestFetchAlwaysTriesOnce(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	var sleeps []time.Duration
	if _, err := newTestFetcher(0, &sleeps).Fetch(context.Background(), srv.URL); err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}

// ctx が切れたらリトライせず即座に諦める。
func TestFetchStopsOnCanceledContext(t *testing.T) {
	var attempts int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var sleeps []time.Duration
	if _, err := newTestFetcher(5, &sleeps).Fetch(ctx, srv.URL); err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if attempts != 0 {
		t.Errorf("attempts = %d, want 0", attempts)
	}
	if len(sleeps) != 0 {
		t.Errorf("sleep called %d times, want 0", len(sleeps))
	}
}
