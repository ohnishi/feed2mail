package main

import (
	"os"
	"path/filepath"
	"time"
)

var subscriptions = []subscription{
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "くるねこ大和", URL: "https://blog.goo.ne.jp/kuru0214/rss2.xml", Fetched: time.Time{}},
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "ITmedia", URL: "https://rss.itmedia.co.jp/rss/2.0/topstory.xml", Fetched: time.Time{}},
	{Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf", Fetched: time.Time{}},
	{Title: "CoinPost", URL: "https://coinpost.jp/?feed=rss2", Fetched: time.Time{}},
	{Title: "コインテレグラフ", URL: "https://jp.cointelegraph.com/rss", Fetched: time.Time{}},
	{Title: "NHKニュース|経済", URL: "https://www3.nhk.or.jp/rss/news/cat5.xml", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "https://randomwalker.blog.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "後藤達也", URL: "https://note.com/goto_finance/rss", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {
	filePath := filepath.Join(dest, "fetchinfo.jsonl")

	if !checkResetModTimeSince(filePath) {
		return nil
	}

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

func checkResetModTimeSince(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}

	nowHour := time.Now().Hour()
	modHour := info.ModTime().Hour()
	// fmt.Println("aaaaaaa", time.Since(info.ModTime()).Hours())
	if nowHour >= 19 && (modHour < 19 || modHour > nowHour) {
		return true
	}

	if time.Since(info.ModTime()).Hours() >= 24 {
		return true
	}

	return false
}
