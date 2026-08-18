package render

import "testing"

func TestBody(t *testing.T) {
	tests := []struct {
		name     string
		sections []Section
		failed   []string
		want     string
	}{
		{
			name: "新着も失敗も無ければ空文字列",
		},
		{
			name: "単一セクション",
			sections: []Section{{
				Title: "Zenn",
				Items: []Item{{Title: "記事A", Link: "https://example.com/a"}},
			}},
			want: "Zenn\n  - 記事A\n    - https://example.com/a",
		},
		{
			name: "複数セクションは空行で区切る",
			sections: []Section{
				{Title: "Zenn", Items: []Item{{Title: "記事A", Link: "https://example.com/a"}}},
				{Title: "はてブ", Items: []Item{{Title: "記事B", Link: "https://example.com/b"}}},
			},
			want: "Zenn\n  - 記事A\n    - https://example.com/a\n\nはてブ\n  - 記事B\n    - https://example.com/b",
		},
		{
			name:     "空のセクションは出力しない",
			sections: []Section{{Title: "Zenn"}, {Title: "はてブ", Items: []Item{{Title: "記事B", Link: "https://example.com/b"}}}},
			want:     "はてブ\n  - 記事B\n    - https://example.com/b",
		},
		{
			name:   "失敗のみでも本文を作る",
			failed: []string{"Zenn", "はてブ"},
			want:   "取得に失敗したフィード:\n  - Zenn\n  - はてブ",
		},
		{
			name:     "新着と失敗を併記する",
			sections: []Section{{Title: "Zenn", Items: []Item{{Title: "記事A", Link: "https://example.com/a"}}}},
			failed:   []string{"はてブ"},
			want:     "Zenn\n  - 記事A\n    - https://example.com/a\n\n取得に失敗したフィード:\n  - はてブ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Body(tt.sections, tt.failed)
			if got != tt.want {
				t.Errorf("Body() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}
