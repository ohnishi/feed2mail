package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadSubscriptionsMissingFile(t *testing.T) {
	got, err := LoadSubscriptions(filepath.Join(t.TempDir(), "absent.jsonl"))
	if err != nil {
		t.Fatalf("LoadSubscriptions() error = %v, want nil", err)
	}
	if got != nil {
		t.Errorf("LoadSubscriptions() = %v, want nil", got)
	}
}

func TestSaveAndLoadSubscriptionsRoundTrip(t *testing.T) {
	// 親ディレクトリが存在しないパスでも書き込めることを併せて確認する。
	path := filepath.Join(t.TempDir(), "cache", "fetched.jsonl")

	want := []Subscription{
		{Title: "Zenn", URL: "https://zenn.dev/feed", Fetched: time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)},
		{Title: "はてブ", URL: "https://b.hatena.ne.jp/hotentry.rss", Fetched: time.Date(2026, 8, 2, 9, 30, 0, 0, time.UTC)},
	}

	if err := SaveSubscriptions(path, want); err != nil {
		t.Fatalf("SaveSubscriptions() error = %v", err)
	}

	got, err := LoadSubscriptions(path)
	if err != nil {
		t.Fatalf("LoadSubscriptions() error = %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("LoadSubscriptions() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Title != want[i].Title || got[i].URL != want[i].URL || !got[i].Fetched.Equal(want[i].Fetched) {
			t.Errorf("subscription[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// 1 レコード 1 行の JSONL であることは git の差分の読みやすさに直結するので固定しておく。
func TestSaveSubscriptionsWritesOneRecordPerLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fetched.jsonl")
	subs := []Subscription{
		{Title: "a", URL: "https://a.example.com"},
		{Title: "b", URL: "https://b.example.com"},
		{Title: "c", URL: "https://c.example.com"},
	}
	if err := SaveSubscriptions(path, subs); err != nil {
		t.Fatalf("SaveSubscriptions() error = %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if n := countNewlines(string(raw)); n != len(subs) {
		t.Errorf("newline count = %d, want %d (content: %q)", n, len(subs), raw)
	}
}

func countNewlines(s string) int {
	n := 0
	for _, r := range s {
		if r == '\n' {
			n++
		}
	}
	return n
}

// 空リストで既存ファイルを消してしまわないことを確認する。
func TestSaveSubscriptionsEmptyKeepsExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fetched.jsonl")
	original := []Subscription{{Title: "a", URL: "https://a.example.com"}}
	if err := SaveSubscriptions(path, original); err != nil {
		t.Fatalf("SaveSubscriptions() error = %v", err)
	}
	if err := SaveSubscriptions(path, nil); err != nil {
		t.Fatalf("SaveSubscriptions(nil) error = %v", err)
	}

	got, err := LoadSubscriptions(path)
	if err != nil {
		t.Fatalf("LoadSubscriptions() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("LoadSubscriptions() len = %d, want 1", len(got))
	}
}
