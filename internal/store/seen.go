package store

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// SeenLink は一度通知対象として評価済みのリンクです。
//
// Source は、そのリンクを最初に見つけた購読フィードの URL です。
// 「このフィードの既読記録がまだ無い＝初回」の判定に使います。
type SeenLink struct {
	Link      string    `json:"link"`
	Source    string    `json:"source"`
	Published time.Time `json:"published"`
	FirstSeen time.Time `json:"first_seen"`
}

// SeenStore は既読リンクの集合です。
//
// 公開時刻の前後だけで新着を判定すると、未来日付や後から遡って追加された記事を
// 取りこぼします。リンク単位で既読を持つことでその両方に耐えられますが、
// 無制限には持てないので一定期間で捨てます。
type SeenStore struct {
	byLink  map[string]SeenLink
	sources map[string]struct{}
}

// NewSeenStore は空の SeenStore を返します。
func NewSeenStore() *SeenStore {
	return &SeenStore{
		byLink:  make(map[string]SeenLink),
		sources: make(map[string]struct{}),
	}
}

// LoadSeen は JSONL 形式の既読リンクを読み込みます。
// ファイルが存在しない場合は空のストアを返します。
func LoadSeen(path string) (*SeenStore, error) {
	s := NewSeenStore()

	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	for d.More() {
		var link SeenLink
		if err := d.Decode(&link); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal seen link: %s", path)
		}
		s.put(link)
	}
	return s, nil
}

func (s *SeenStore) put(link SeenLink) {
	s.byLink[link.Link] = link
	if link.Source != "" {
		s.sources[link.Source] = struct{}{}
	}
}

// Has はリンクが既読かどうかを返します。
func (s *SeenStore) Has(link string) bool {
	_, ok := s.byLink[link]
	return ok
}

// HasSource は指定フィード由来の既読記録が 1 件でもあるかを返します。
// false の場合、そのフィードはまだ既読リンクを持たない（初回）とみなせます。
func (s *SeenStore) HasSource(source string) bool {
	_, ok := s.sources[source]
	return ok
}

// Add はリンクを既読として記録します。既に記録済みなら何もしません。
func (s *SeenStore) Add(link, source string, published, now time.Time) {
	if s.Has(link) {
		return
	}
	s.put(SeenLink{Link: link, Source: source, Published: published, FirstSeen: now})
}

// Prune は cutoff より前に記録されたリンクを捨て、捨てた件数を返します。
func (s *SeenStore) Prune(cutoff time.Time) int {
	var pruned int
	for key, link := range s.byLink {
		if link.FirstSeen.Before(cutoff) {
			delete(s.byLink, key)
			pruned++
		}
	}
	if pruned > 0 {
		s.rebuildSources()
	}
	return pruned
}

func (s *SeenStore) rebuildSources() {
	s.sources = make(map[string]struct{}, len(s.sources))
	for _, link := range s.byLink {
		if link.Source != "" {
			s.sources[link.Source] = struct{}{}
		}
	}
}

// Len は保持している既読リンクの件数を返します。
func (s *SeenStore) Len() int { return len(s.byLink) }

// All は記録順（FirstSeen 昇順、同時刻ならリンク昇順）に既読リンクを返します。
func (s *SeenStore) All() []SeenLink {
	links := make([]SeenLink, 0, len(s.byLink))
	for _, link := range s.byLink {
		links = append(links, link)
	}
	// 新しい記録がファイル末尾に寄るよう並べる。git の差分が追記だけで済む。
	sort.Slice(links, func(i, j int) bool {
		if !links[i].FirstSeen.Equal(links[j].FirstSeen) {
			return links[i].FirstSeen.Before(links[j].FirstSeen)
		}
		return links[i].Link < links[j].Link
	})
	return links
}

// SaveSeen は既読リンクを JSONL 形式で書き出します。
// 空でも書き出します。ファイルの存在自体が「初回ではない」ことの印になるためです。
func SaveSeen(path string, s *SeenStore) error {
	return writeJSONL(path, func(enc *json.Encoder) error {
		for _, link := range s.All() {
			if err := enc.Encode(link); err != nil {
				return errors.Wrapf(err, "failed to write seen link: %s", link.Link)
			}
		}
		return nil
	})
}
