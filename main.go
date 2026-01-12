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
	// {Title: "CNET Japan", URL: "https://feeds.japan.cnet.com/rss/cnet/all.rdf", Fetched: time.Time{}},
	// {Title: "INTERNET Watch", URL: "https://internet.watch.impress.co.jp/data/rss/1.0/iw/feed.rdf", Fetched: time.Time{}},
	// {Title: "GIZMODO JAPAN", URL: "https://www.gizmodo.jp/index.xml", Fetched: time.Time{}},
	// {Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf", Fetched: time.Time{}},
	// {Title: "Qiita", URL: "https://qiita.com/popular-items/feed.atom", Fetched: time.Time{}},
	{Title: "Zenn", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "市況かぶ全力２階建", URL: "https://kabumatome.doorblog.jp/atom.xml", Fetched: time.Time{}},
	{Title: "勝間和代が徹底的にマニアックな話をアップするブログ", URL: "https://katsumakazuyo.hatenablog.com/rss", Fetched: time.Time{}},
	{Title: "The GitHub Blog", URL: "https://github.blog/feed/", Fetched: time.Time{}},

	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "チェンソーマン", URL: "https://shonenjumpplus.com/rss/series/3270375685341574005", Fetched: time.Time{}},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490", Fetched: time.Time{}},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745", Fetched: time.Time{}},

	{Title: "しろまるnote", URL: "https://note.com/asuka_shiromaru/rss", Fetched: time.Time{}},
	// {Title: "セールモンスター", URL: "https://note.com/_salemonster/rss", Fetched: time.Time{}},
	{Title: "久松剛", URL: "https://note.com/makaibito/rss", Fetched: time.Time{}},
	{Title: "わさびん土屋", URL: "https://note.com/wasabinbin/rss", Fetched: time.Time{}},
	{Title: "ニケちゃん", URL: "https://note.com/nike_cha_n/rss", Fetched: time.Time{}},
	{Title: "安野貴博", URL: "https://note.com/takahiroanno/rss", Fetched: time.Time{}},
	{Title: "牛尾剛", URL: "https://note.com/simplearchitect/rss", Fetched: time.Time{}},
	{Title: "入江慎吾", URL: "https://note.com/iritec/rss", Fetched: time.Time{}},
	{Title: "のすけ", URL: "https://note.com/nosuke0926/rss", Fetched: time.Time{}},
	{Title: "五月（片山晃）", URL: "https://note.com/may5x/rss", Fetched: time.Time{}},
	{Title: "唐鎌大輔", URL: "https://note.com/dkarakama/rss", Fetched: time.Time{}},

	{Title: "オライリー", URL: "https://www.oreilly.co.jp/catalog/soon.xml", Fetched: time.Time{}},
	{Title: "技術評論社", URL: "https://gihyo.jp/book/feed/rss1", Fetched: time.Time{}},

	{Title: "Yahoo!ニュース", URL: "https://news.yahoo.co.jp/rss/topics/top-picks.xml", Fetched: time.Time{}},
	// {Title: "日経電子版", URL: "https://news.google.com/rss/search?q=site:nikkei.com&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},
	{Title: "NewsPicks 記事", URL: "https://news.google.com/rss/search?q=site:newspicks.com&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},
	{Title: "四季報オンライン", URL: "https://news.google.com/rss/search?q=site:shikiho.toyokeizai.net&hl=ja&gl=JP&ceid=JP:ja", Fetched: time.Time{}},

	{Title: "楽待", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPMJKbrxtpARoTd1b49iUvA", Fetched: time.Time{}},
	{Title: "KabuBerry Channel", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC3z7uU4ybkh9Pkz5A0-QGYA", Fetched: time.Time{}},
	{Title: "湘南投資勉強会オンライン", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCLJ7mt8d0GbJs74PTPlREYg", Fetched: time.Time{}},
	{Title: "不動産Gメン滝島", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCv9GARGGn0LpNq7QHQVZp2A", Fetched: time.Time{}},
	{Title: "田端大学", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC7sEB_ylMuHJD4TjF4Ag1nw", Fetched: time.Time{}},
	{Title: "中国まる見え情報局", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgv73JUFDUd_nIulI_VECKQ", Fetched: time.Time{}},
	{Title: "勝間和代が徹底的にマニアックな話をするYouTube", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCWoiNwdr7EEjgs2waxe_QpA", Fetched: time.Time{}},
	{Title: "テレ東AIアカデミー【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCtadcNeydAUCkcUF5t5YNtw", Fetched: time.Time{}},
	{Title: "ReHacQ−リハック−【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCG_oqDSlIYEspNpd2H4zWhw", Fetched: time.Time{}},
	{Title: "テレ東BIZ ダイジェスト", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCkKVQ_GNjd8FbAuT6xDcWgg", Fetched: time.Time{}},
	{Title: "PIVOT 公式チャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC8yHePe_RgUBE-waRWy6olw", Fetched: time.Time{}},
	{Title: "文藝春秋PLUS 公式チャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCTcJVDgfw411rqWCHb-16lg", Fetched: time.Time{}},
	{Title: "TBS NEWS DIG Powered by JNN", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC6AG81pAkf6Lbi_1VC5NmPA", Fetched: time.Time{}},
	{Title: "ダイヤモンド社 THE BOOKS", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UChyLk7utaMRsA7qZ6zT8UUg", Fetched: time.Time{}},
	{Title: "NewsPicks 動画", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfTnJmRQP79C4y_BMF_XrlA", Fetched: time.Time{}},

	{Title: "物販NAVI", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMO-lX7Y0K7c-BaV9RatboQ", Fetched: time.Time{}},
	{Title: "中古カメラ現物投資家 船田ひろし", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcph2C76gdSHXUbd3wooNfQ", Fetched: time.Time{}},
	{Title: "脱サラした男のネットショップな日常", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCfCkw1ugUhc9Nw_wu1cKHpw", Fetched: time.Time{}},
	{Title: "かわしま＠中国輸入", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCW_Ey8M8I3o1JCB2OyCYzHA", Fetched: time.Time{}},
	{Title: "こーいち【中国輸入Amazon販売】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCqmEmvM6rt3iTYiTv-PRI0g", Fetched: time.Time{}},
	{Title: "Tasチャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC96v4EM1wETxpqrK2WqrFNw", Fetched: time.Time{}},
	{Title: "Tasの物販作業部屋", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCmzbvkNJEubCV9pKfYfPlpA", Fetched: time.Time{}},
	{Title: "TasのEC支援", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCtFFMu0IO0CMmHkZyOqTMjg", Fetched: time.Time{}},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA", Fetched: time.Time{}},
	{Title: "Amazonで売る【公式】", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},

	{Title: "MLBコアラ", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCQTnU6pFNEsOJWHSdj7Y5qA", Fetched: time.Time{}},
	{Title: "全力野球少年", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMH6sSHDQt6yDwFjbask4OQ", Fetched: time.Time{}},
	{Title: "ミスターフルスイングch", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCoryrORUTuHoSkryhk0K8bA", Fetched: time.Time{}},
	{Title: "プロ野球のおもいでch", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCArbMa_kKr91pH0KkzQsEkA", Fetched: time.Time{}},

	{Title: "TECH WORLD", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCISDrqLMNq3w9AZ4otdoRuA", Fetched: time.Time{}},
	{Title: "テックナビ", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCB3mOPzHaD-K1n7YlubgOqA", Fetched: time.Time{}},
	{Title: "中島聡のLife is Beautiful", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCtjRA-7EuBmyWMyty1HZAPQ", Fetched: time.Time{}},
	{Title: "中村祐太のFindUアカデミー", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCNyN5h_EJ1ktJoTODTD6DyA", Fetched: time.Time{}},
	{Title: "入江慎吾", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC3wNxV5ZWOL9kf4v2oXGQbw", Fetched: time.Time{}},
	{Title: "AI Engineer", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCLKPca3kwwd-B59HNr-_lvA", Fetched: time.Time{}},
	{Title: "Vanessa Wingårdh", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC2Bm4qXVUx9Md-hmHz3WuXw", Fetched: time.Time{}},

	{Title: "Naokiman Show", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC4lN5sizuJraSHqy99xTy6Q", Fetched: time.Time{}},
	{Title: "Naokiman 2nd Channel", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCT9cFrlnfy1-XLf0Vuy0Zzw", Fetched: time.Time{}},

	{Title: "ゆうひなたチャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC-MjY-iABWP-JZOzwnuwLWg", Fetched: time.Time{}},
	{Title: "Bappa Shota", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcrpFRRYkH185Xb8D-JQT7Q", Fetched: time.Time{}},
	{Title: "ビビりの家族が行くキャンピングカー生活", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC-YMyraUiYnA1c8IMy-JjGA", Fetched: time.Time{}},
	{Title: "スタイリスト大山シュンのメンズ服講座", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCRf_7nS9lxLZFfK6kJBvK5A", Fetched: time.Time{}},
	{Title: "Living Big In A Tiny House", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCoNTMWgGuXtGPLv9UeJZwBw", Fetched: time.Time{}},
	{Title: "JOJO Channel<", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCx8E8jP6tH1ZSRVYklgu_Dg", Fetched: time.Time{}},
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
