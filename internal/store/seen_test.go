package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var base = time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)

func TestLoadSeenMissingFileReturnsEmptyStore(t *testing.T) {
	s, err := LoadSeen(filepath.Join(t.TempDir(), "absent.jsonl"))
	if err != nil {
		t.Fatalf("LoadSeen() error = %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("Len() = %d, want 0", s.Len())
	}
}

func TestAddIsIdempotent(t *testing.T) {
	s := NewSeenStore()
	s.Add("https://example.com/a", "https://src.example.com", base, base)
	s.Add("https://example.com/a", "https://other.example.com", base.Add(time.Hour), base.Add(time.Hour))

	if s.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", s.Len())
	}
	// 最初の記録が残る。FirstSeen が上書きされると保持期間が伸び続けてしまう。
	if got := s.All()[0]; got.Source != "https://src.example.com" || !got.FirstSeen.Equal(base) {
		t.Errorf("record = %+v, want the first one", got)
	}
}

func TestHasSource(t *testing.T) {
	s := NewSeenStore()
	if s.HasSource("https://src.example.com") {
		t.Error("HasSource() = true on empty store, want false")
	}
	s.Add("https://example.com/a", "https://src.example.com", base, base)
	if !s.HasSource("https://src.example.com") {
		t.Error("HasSource() = false after Add, want true")
	}
	if s.HasSource("https://another.example.com") {
		t.Error("HasSource() = true for unknown source, want false")
	}
}

func TestPruneDropsOldRecordsAndSources(t *testing.T) {
	s := NewSeenStore()
	s.Add("https://example.com/old", "https://old.example.com", base, base.Add(-40*24*time.Hour))
	s.Add("https://example.com/new", "https://new.example.com", base, base)

	if pruned := s.Prune(base.Add(-30 * 24 * time.Hour)); pruned != 1 {
		t.Fatalf("Prune() = %d, want 1", pruned)
	}
	if s.Has("https://example.com/old") {
		t.Error("old record survived Prune()")
	}
	if !s.Has("https://example.com/new") {
		t.Error("recent record was pruned")
	}
	// Source の索引も追随しないと、初回判定が誤る。
	if s.HasSource("https://old.example.com") {
		t.Error("HasSource() = true for pruned source, want false")
	}
	if !s.HasSource("https://new.example.com") {
		t.Error("HasSource() = false for surviving source, want true")
	}
}

func TestSaveAndLoadSeenRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache", "seen.jsonl")
	s := NewSeenStore()
	s.Add("https://example.com/a", "https://src.example.com", base.Add(-time.Hour), base)
	s.Add("https://example.com/b", "https://src.example.com", base, base)

	if err := SaveSeen(path, s); err != nil {
		t.Fatalf("SaveSeen() error = %v", err)
	}

	got, err := LoadSeen(path)
	if err != nil {
		t.Fatalf("LoadSeen() error = %v", err)
	}
	if got.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", got.Len())
	}
	if !got.Has("https://example.com/a") || !got.Has("https://example.com/b") {
		t.Error("round trip lost a link")
	}
	if !got.HasSource("https://src.example.com") {
		t.Error("round trip lost the source index")
	}
}

// 新しい記録が末尾に寄れば、git の差分は追記だけで済む。
func TestAllSortsByFirstSeenThenLink(t *testing.T) {
	s := NewSeenStore()
	s.Add("https://example.com/z", "src", base, base.Add(-2*time.Hour))
	s.Add("https://example.com/b", "src", base, base)
	s.Add("https://example.com/a", "src", base, base)

	want := []string{"https://example.com/z", "https://example.com/a", "https://example.com/b"}
	got := s.All()
	for i := range want {
		if got[i].Link != want[i] {
			t.Errorf("All()[%d].Link = %q, want %q", i, got[i].Link, want[i])
		}
	}
}

// 空でもファイルを作る。存在自体が「初回ではない」ことの印になる。
func TestSaveSeenWritesFileEvenWhenEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seen.jsonl")
	if err := SaveSeen(path, NewSeenStore()); err != nil {
		t.Fatalf("SaveSeen() error = %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if strings.TrimSpace(string(raw)) != "" {
		t.Errorf("file content = %q, want empty", raw)
	}
}
