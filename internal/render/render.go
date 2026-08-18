// Package render は通知メールの本文を組み立てます。
// 副作用を持たない純粋な変換に閉じており、出力はテストで固定しています。
package render

import "strings"

// Item は本文に載せる 1 記事です。
type Item struct {
	Title string
	Link  string
}

// Section はフィード 1 件分の新着記事です。
type Section struct {
	Title string
	Items []Item
}

// failedHeading は取得に失敗したフィードを並べる見出しです。
const failedHeading = "取得に失敗したフィード:"

// Body は新着セクションと失敗フィード一覧からメール本文を組み立てます。
// 新着も失敗も無い場合は空文字列を返します。呼び出し側はこれを送信スキップの判断に使えます。
func Body(sections []Section, failed []string) string {
	var blocks []string

	for _, section := range sections {
		if len(section.Items) == 0 {
			continue
		}
		lines := make([]string, 0, len(section.Items)*2+1)
		lines = append(lines, section.Title)
		for _, item := range section.Items {
			lines = append(lines, "  - "+item.Title)
			lines = append(lines, "    - "+item.Link)
		}
		blocks = append(blocks, strings.Join(lines, "\n"))
	}

	if len(failed) > 0 {
		lines := make([]string, 0, len(failed)+1)
		lines = append(lines, failedHeading)
		for _, title := range failed {
			lines = append(lines, "  - "+title)
		}
		blocks = append(blocks, strings.Join(lines, "\n"))
	}

	return strings.Join(blocks, "\n\n")
}
