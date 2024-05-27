package main

import (
	"os"
	"path/filepath"
	"time"
)

var subscriptions = []subscription{
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "くるねこ大和", URL: "https://blog.goo.ne.jp/kuru0214/rss2.xml", Fetched: time.Time{}},
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "https://randomwalker.blog.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "CiLEL", URL: "https://cilel.jp/feed/", Fetched: time.Time{}},
	{Title: "ログミーTech", URL: "https://logmi.jp/feed/public-tech.xml", Fetched: time.Time{}},
	{Title: "ガジェラン", URL: "https://gadgelaun.com/?feed=rss2", Fetched: time.Time{}},
	{Title: "IIJ Engineers Blog", URL: "https://eng-blog.iij.ad.jp/feed", Fetched: time.Time{}},
	{Title: "Money Forward Developers Blog", URL: "https://moneyforward-dev.jp/feed", Fetched: time.Time{}},
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "ITmedia", URL: "https://rss.itmedia.co.jp/rss/2.0/topstory.xml", Fetched: time.Time{}},
	{Title: "価格.comマガジン", URL: "https://kakakumag.com/rss/", Fetched: time.Time{}},
	{Title: "ダイヤモンド・オンライン", URL: "https://diamond.jp/list/feed/rss/dol", Fetched: time.Time{}},
	{Title: "ザイ・オンライン", URL: "https://diamond.jp/zai/list/feed/rsszol", Fetched: time.Time{}},
	{Title: "東洋経済オンライン", URL: "https://toyokeizai.net/list/feed/rss", Fetched: time.Time{}},
	{Title: "鈴木正行", URL: "https://note.com/redturtle0721/rss", Fetched: time.Time{}},
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
		return false
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
