package main

import (
	"path/filepath"
	"time"

	"github.com/ohnishi/feed/models"
)

var subscriptions = []models.Subscription{
	{Title: "はてブ（総合）", URL: "https://b.hatena.ne.jp/hotentry.rss", Fetched: time.Time{}},
	{Title: "はてブ（テクノロジー）", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},

	{Title: "Zenn", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},

	{Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf", Fetched: time.Time{}},
	{Title: "AI Watch", URL: "https://ai.watch.impress.co.jp/data/rss/1.0/aiw/feed.rdf", Fetched: time.Time{}},

	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "くるねこ大和", URL: "https://rssblog.ameba.jp/kuru0214neko/rss20.xml", Fetched: time.Time{}},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490", Fetched: time.Time{}},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745", Fetched: time.Time{}},

	{Title: "物販NAVI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMO-lX7Y0K7c-BaV9RatboQ", Fetched: time.Time{}},
	{Title: "中古カメラ現物投資家 船田ひろし", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcph2C76gdSHXUbd3wooNfQ", Fetched: time.Time{}},
	{Title: "脱サラした男のネットショップな日常", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfCkw1ugUhc9Nw_wu1cKHpw", Fetched: time.Time{}},
	{Title: "かわしま＠中国輸入", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCW_Ey8M8I3o1JCB2OyCYzHA", Fetched: time.Time{}},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA", Fetched: time.Time{}},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {
	filePath := filepath.Join(dest)

	oldSubscriptions, err := models.ReadSubscriptions(filePath)
	if err != nil {
		return err
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// tzdata を持たない環境へのフォールバック
		jst = time.FixedZone("JST", 9*60*60)
	}

	var newSubscriptions []models.Subscription
	for _, s := range subscriptions {
		isExists := false
		for _, old := range oldSubscriptions {
			if s.URL == old.URL {
				// 取得済み時刻だけを引き継ぎ、Title は購読リスト側の定義を優先する
				s.Fetched = old.Fetched
				isExists = true
				break
			}
		}
		if !isExists {
			s.Fetched = time.Now().In(jst).AddDate(0, 0, -1)
		}
		newSubscriptions = append(newSubscriptions, s)
	}

	return models.WriteSubscription(filePath, newSubscriptions)
}

const (
	maxRetry = 5
	saveFile = "cache/fetched.jsonl"
)

func main() {
	err := resetSubscriptions(saveFile)
	if err != nil {
		panic(err)
	}

	err = models.FeedToMail(saveFile, maxRetry)
	if err != nil {
		panic(err)
	}
}
