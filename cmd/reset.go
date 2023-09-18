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
	{Title: "Yahoo!ニュース・トピックス - IT", URL: "https://news.yahoo.co.jp/rss/topics/it.xml", Fetched: time.Time{}},
	{Title: "Yahoo!ニュース・トピックス - 経済", URL: "https://news.yahoo.co.jp/rss/topics/business.xml", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "そばに", URL: "https://sobani.co.jp/feed", Fetched: time.Time{}},
	{Title: "CiLEL", URL: "https://cilel.jp/feed/", Fetched: time.Time{}},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},
	{Title: "Rebuild", URL: "https://feeds.rebuild.fm/rebuildfm", Fetched: time.Time{}},
	{Title: "入江 慎吾", URL: "https://iritec.jp/feed/", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {

	return writeSubscription(filepath.Join(dest, "fetchinfo.json"), subscriptions)
}
