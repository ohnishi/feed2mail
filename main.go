package main

import (
	"log"
	"time"

	"github.com/ohnishi/feed/models"
)

var subscriptions = []models.Subscription{
	{Title: "はてブ（総合）", URL: "https://b.hatena.ne.jp/hotentry.rss"},
	{Title: "はてブ（テクノロジー）", URL: "https://b.hatena.ne.jp/hotentry/it.rss"},

	{Title: "Zenn", URL: "https://zenn.dev/feed"},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml"},

	{Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf"},
	{Title: "AI Watch", URL: "https://ai.watch.impress.co.jp/data/rss/1.0/aiw/feed.rdf"},

	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf"},
	{Title: "くるねこ大和", URL: "https://rssblog.ameba.jp/kuru0214neko/rss20.xml"},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490"},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745"},

	{Title: "物販NAVI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMO-lX7Y0K7c-BaV9RatboQ"},
	{Title: "中古カメラ現物投資家 船田ひろし", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcph2C76gdSHXUbd3wooNfQ"},
	{Title: "脱サラした男のネットショップな日常", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfCkw1ugUhc9Nw_wu1cKHpw"},
	{Title: "かわしま＠中国輸入", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCW_Ey8M8I3o1JCB2OyCYzHA"},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA"},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA"},
}

func resetSubscriptions(filePath string) error {
	oldSubscriptions, err := models.ReadSubscriptions(filePath)
	if err != nil {
		return err
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// tzdata を持たない環境へのフォールバック
		jst = time.FixedZone("JST", 9*60*60)
	}

	fetchedByURL := make(map[string]time.Time, len(oldSubscriptions))
	for _, old := range oldSubscriptions {
		fetchedByURL[old.URL] = old.Fetched
	}

	newSubscriptions := make([]models.Subscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		if fetched, ok := fetchedByURL[s.URL]; ok {
			// 取得済み時刻だけを引き継ぎ、Title は購読リスト側の定義を優先する
			s.Fetched = fetched
		} else {
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
	if err := resetSubscriptions(saveFile); err != nil {
		log.Fatalf("failed to reset subscriptions: %v", err)
	}

	if err := models.FeedToMail(saveFile, maxRetry); err != nil {
		log.Fatalf("failed to deliver feed: %v", err)
	}
}
