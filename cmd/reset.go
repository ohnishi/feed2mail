package main

import (
	"path/filepath"
	"time"
)

var subscriptions = []subscription{
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "くるねこ大和", URL: "https://blog.goo.ne.jp/kuru0214/rss2.xml", Fetched: time.Time{}},
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "https://randomwalker.blog.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "CiLEL", URL: "https://cilel.jp/feed/", Fetched: time.Time{}},
	{Title: "独立を楽しくするブログ", URL: "https://www.ex-it-blog.com/feed/", Fetched: time.Time{}},
	{Title: "PCまなぶ", URL: "https://pcmanabu.com/feed/", Fetched: time.Time{}},
	{Title: "ログミーTech", URL: "https://logmi.jp/feed/public-tech.xml", Fetched: time.Time{}},
	{Title: "ガジェラン", URL: "https://gadgelaun.com/?feed=rss2", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {
	filePath := filepath.Join(dest, "fetchinfo.jsonl")
	oldSubscriptions, err := readSubscriptions(filePath)
	if err != nil {
		return err
	}

	var newSubscriptions []subscription
	for _, s := range subscriptions {
		isExists := false
		for _, old := range oldSubscriptions {
			if s.URL == old.URL {
				newSubscriptions = append(newSubscriptions, old)
				isExists = true
				break
			}
		}
		if !isExists {
			newSubscriptions = append(newSubscriptions, s)
		}
	}

	return writeSubscription(filePath, newSubscriptions)
}
