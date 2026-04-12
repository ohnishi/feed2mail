package main

import (
	"path/filepath"
	"time"

	"github.com/ohnishi/feed/models"
)

var subscriptions = []models.Subscription{
	{Title: "はてブ（総合）", URL: "https://b.hatena.ne.jp/hotentry.rss", Fetched: time.Time{}},
	{Title: "はてブ（テクノロジー）", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "ITmedia", URL: "https://rss.itmedia.co.jp/rss/2.0/topstory.xml", Fetched: time.Time{}},
	{Title: "Zenn", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	// {Title: "市況かぶ全力２階建", URL: "https://kabumatome.doorblog.jp/atom.xml", Fetched: time.Time{}},
	// {Title: "勝間和代が徹底的にマニアックな話をアップするブログ", URL: "https://katsumakazuyo.hatenablog.com/rss", Fetched: time.Time{}},
	{Title: "The GitHub Blog", URL: "https://github.blog/feed/", Fetched: time.Time{}},

	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490", Fetched: time.Time{}},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745", Fetched: time.Time{}},

	{Title: "しろまるnote", URL: "https://note.com/asuka_shiromaru/rss", Fetched: time.Time{}},
	{Title: "久松剛", URL: "https://note.com/makaibito/rss", Fetched: time.Time{}},
	{Title: "わさびん土屋", URL: "https://note.com/wasabinbin/rss", Fetched: time.Time{}},
	{Title: "ニケちゃん", URL: "https://note.com/nike_cha_n/rss", Fetched: time.Time{}},
	{Title: "安野貴博", URL: "https://note.com/takahiroanno/rss", Fetched: time.Time{}},
	{Title: "牛尾剛", URL: "https://note.com/simplearchitect/rss", Fetched: time.Time{}},
	{Title: "入江慎吾", URL: "https://note.com/iritec/rss", Fetched: time.Time{}},
	{Title: "のすけ", URL: "https://note.com/nosuke0926/rss", Fetched: time.Time{}},
	{Title: "五月（片山晃）", URL: "https://note.com/may5x/rss", Fetched: time.Time{}},
	{Title: "唐鎌大輔", URL: "https://note.com/dkarakama/rss", Fetched: time.Time{}},
	{Title: "ジェネトピ", URL: "https://note.com/genai_topic/rss", Fetched: time.Time{}},

	{Title: "オライリー", URL: "https://www.oreilly.co.jp/catalog/soon.xml", Fetched: time.Time{}},
	{Title: "技術評論社", URL: "https://gihyo.jp/book/feed/rss1", Fetched: time.Time{}},

	{Title: "Yahoo!ニュース", URL: "https://news.yahoo.co.jp/rss/topics/top-picks.xml", Fetched: time.Time{}},
	// {Title: "日経電子版", URL: "https://news.google.com/rss/search?q=site:nikkei.com&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},
	// {Title: "NewsPicks 記事", URL: "https://news.google.com/rss/search?q=site:newspicks.com&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},
	// {Title: "四季報オンライン", URL: "https://news.google.com/rss/search?q=site:shikiho.toyokeizai.net&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},

	{Title: "物販NAVI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMO-lX7Y0K7c-BaV9RatboQ", Fetched: time.Time{}},
	{Title: "中古カメラ現物投資家 船田ひろし", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcph2C76gdSHXUbd3wooNfQ", Fetched: time.Time{}},
	{Title: "脱サラした男のネットショップな日常", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfCkw1ugUhc9Nw_wu1cKHpw", Fetched: time.Time{}},
	{Title: "かわしま＠中国輸入", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCW_Ey8M8I3o1JCB2OyCYzHA", Fetched: time.Time{}},
	{Title: "こーいち【中国輸入Amazon販売】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCqmEmvM6rt3iTYiTv-PRI0g", Fetched: time.Time{}},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA", Fetched: time.Time{}},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},
	{Title: "ERESA イーリサ ®公式チャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCiiO8UnzTfwk1EflGir1VyQ", Fetched: time.Time{}},

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
