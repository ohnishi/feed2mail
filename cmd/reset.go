package main

import (
	"path/filepath"
	"time"
)

var subscriptions = []subscription{
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "http://randomwalker.blog19.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "はてなブックマーク - 人気エントリー - テクノロジー", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {

	return writeSubscription(filepath.Join(dest, "fetchinfo.json"), subscriptions)
}
