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
	{Title: "Rebuild Podcast", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCty73f38mfBICYNZAP839fA", Fetched: time.Time{}},
	{Title: "ムーザルちゃんねる", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCLPHXwLp90A5R69Eltxo-sg", Fetched: time.Time{}},

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
	{Title: "こーいち【中国輸入Amazon販売】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCqmEmvM6rt3iTYiTv-PRI0g", Fetched: time.Time{}},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA", Fetched: time.Time{}},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},
	{Title: "ERESA イーリサ ®公式チャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCiiO8UnzTfwk1EflGir1VyQ", Fetched: time.Time{}},
	{Title: "ゲーハー田村", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCVWsF_cg3piJ7rRYz7bItbQ", Fetched: time.Time{}},

	{Title: "マーケットマスターズ", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCXvjRTXoDa8tKwdkTaukGug", Fetched: time.Time{}},
	{Title: "stockvoice", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCeW7dO9lXpvUQJcMEP1tv0g", Fetched: time.Time{}},
	{Title: "楽待 RAKUMACHI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPMJKbrxtpARoTd1b49iUvA", Fetched: time.Time{}},
	{Title: "日経CNBC", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UClVsQnfs-jKkjKmUKUHnT2g", Fetched: time.Time{}},
	{Title: "TBS CROSS DIG with Bloomberg", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCeCmAYh1ylwIsgGrmqaklzg", Fetched: time.Time{}},

	{Title: "聖女イトコイザーさん【仮】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC5Ou3Sy_86_pfHlZf7v-REw", Fetched: time.Time{}},

	// {Title: "", URL: "", Fetched: time.Time{}},
}

func resetSubscriptions(dest string) error {
	filePath := filepath.Join(dest)

	oldSubscriptions, err := models.ReadSubscriptions(filePath)
	if err != nil {
		return err
	}

	var newSubscriptions []models.Subscription
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
			jst, _ := time.LoadLocation("Asia/Tokyo")
			s.Fetched = time.Now().In(jst).AddDate(0, 0, -1)
			newSubscriptions = append(newSubscriptions, s)
		}
	}

	return models.WriteSubscription(filePath, newSubscriptions)
}

const (
	maxRetry = 10
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
