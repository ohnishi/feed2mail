package main

import (
	"path/filepath"
	"time"

	"github.com/ohnishi/feed/models"
)

var subscriptions = []models.Subscription{
	{Title: "NHKニュース", URL: "https://www3.nhk.or.jp/rss/news/cat0.xml", Fetched: time.Time{}},
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "ITmedia", URL: "https://rss.itmedia.co.jp/rss/2.0/topstory.xml", Fetched: time.Time{}},
	{Title: "CNET Japan", URL: "https://feeds.japan.cnet.com/rss/cnet/all.rdf", Fetched: time.Time{}},
	{Title: "INTERNET Watch", URL: "https://internet.watch.impress.co.jp/data/rss/1.0/iw/feed.rdf", Fetched: time.Time{}},
	{Title: "PC Watch", URL: "https://pc.watch.impress.co.jp/data/rss/1.0/pcw/feed.rdf", Fetched: time.Time{}},
	{Title: "GIGAZINE", URL: "https://gigazine.net/news/rss_2.0/", Fetched: time.Time{}},
	{Title: "GIZMODO JAPAN", URL: "https://www.gizmodo.jp/index.xml", Fetched: time.Time{}},
	{Title: "窓の杜", URL: "https://forest.watch.impress.co.jp/data/rss/1.0/wf/feed.rdf", Fetched: time.Time{}},
	{Title: "ROOMIE", URL: "https://www.roomie.jp/feed/", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "https://randomwalker.blog.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "チェンソーマン", URL: "https://shonenjumpplus.com/rss/series/3270375685341574005", Fetched: time.Time{}},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490", Fetched: time.Time{}},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745", Fetched: time.Time{}},
	{Title: "しろまるnote", URL: "https://note.com/asuka_shiromaru/rss", Fetched: time.Time{}},
	{Title: "セールモンスター", URL: "https://note.com/_salemonster/rss", Fetched: time.Time{}},
	{Title: "セラースプライト", URL: "https://note.com/sellersprite/rss", Fetched: time.Time{}},
	{Title: "久松剛", URL: "https://note.com/makaibito/rss", Fetched: time.Time{}},
	{Title: "わさびん土屋", URL: "https://note.com/wasabinbin/rss", Fetched: time.Time{}},
	{Title: "ニケちゃん", URL: "https://note.com/nike_cha_n/rss", Fetched: time.Time{}},
	{Title: "安野貴博", URL: "https://note.com/takahiroanno/rss", Fetched: time.Time{}},
	{Title: "牛尾 剛", URL: "https://note.com/simplearchitect/rss", Fetched: time.Time{}},
	{Title: "現海秀哲物販Note", URL: "https://note.com/genkai_hideaki/rss", Fetched: time.Time{}},
	{Title: "スーパーポテサラアニメ", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC1ees-fSpiNAIy7O7idjLuA", Fetched: time.Time{}},
	{Title: "文藝春秋PLUS", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCTcJVDgfw411rqWCHb-16lg", Fetched: time.Time{}},
	{Title: "sexyfitness", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC2b1bVPPJgpf4qb6Asij33A", Fetched: time.Time{}},
	{Title: "Marques Brownlee", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCBJycsmduvYEL83R_U4JriQ", Fetched: time.Time{}},
	{Title: "セールモンスター", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC3Yce3hF7z7RXZ8P6s297RA", Fetched: time.Time{}},
	{Title: "るりチャンネル", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCrQsBzv2Lfq4RrVnugT87Lw", Fetched: time.Time{}},
	{Title: "中島聡", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCtjRA-7EuBmyWMyty1HZAPQ", Fetched: time.Time{}},
	{Title: "楽待", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPMJKbrxtpARoTd1b49iUvA", Fetched: time.Time{}},
	{Title: "メンズ服講座", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCRf_7nS9lxLZFfK6kJBvK5A", Fetched: time.Time{}},
	{Title: "PAD", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC7WpJ8eZESNDtO2uALSjigQ", Fetched: time.Time{}},
	{Title: "ヒロ税理士", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC1nVq5-LTEdUVmr5FLfkd7A", Fetched: time.Time{}},
	{Title: "Googleシートマスターch", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcJ5SvYn_PxgIaqg0oERpAw", Fetched: time.Time{}},
	{Title: "吉田製作所Y", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC9WJo5ZJVXMZiA5XV2jLx5Q", Fetched: time.Time{}},
	{Title: "Seller Sprite", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC7UGrov7pv43Vx75Xqj7ZfQ", Fetched: time.Time{}},
	{Title: "AI駆動開発勉強会", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCTILP8NYy3l5xMb87IkbH3g", Fetched: time.Time{}},
	{Title: "入江慎吾", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC3wNxV5ZWOL9kf4v2oXGQbw", Fetched: time.Time{}},
	{Title: "中村祐太", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCNyN5h_EJ1ktJoTODTD6DyA", Fetched: time.Time{}},
	{Title: "定額制リフォーム", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCSpi9noH1wwFGkpKNFr6J-Q", Fetched: time.Time{}},
	{Title: "Amazonで売る", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCPBWwGeO-vbrsOCGrd0_rrA", Fetched: time.Time{}},
	{Title: "Naokiman Show", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC4lN5sizuJraSHqy99xTy6Q", Fetched: time.Time{}},
	{Title: "オタク会計士ch", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCMAEQdzGckZ9FMWJv8tz2zA", Fetched: time.Time{}},
	{Title: "RAKUMART", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCgW9fypC1UpPsIF1m4Ph5eA", Fetched: time.Time{}},
	{Title: "Newbee", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCS0Nxf1j4Jtfx2i9qTksezg", Fetched: time.Time{}},
	{Title: "吉田研究所", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCLE2PTNDt--NLGLVbfOfHkQ", Fetched: time.Time{}},
	{Title: "こぶしの板さん", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCVDMeXIV3elsm8mjWre-qIA", Fetched: time.Time{}},
	{Title: "こぶしの山で!", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UC5_q8XHadgecavUCCYm96Gw", Fetched: time.Time{}},
	{Title: "田舎そば川原", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCDzOfd9304g0cWmSameJNWw", Fetched: time.Time{}},
	{Title: "Naokiman 2nd Channel", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCT9cFrlnfy1-XLf0Vuy0Zzw", Fetched: time.Time{}},
	{Title: "Bappa Shota", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCcrpFRRYkH185Xb8D-JQT7Q", Fetched: time.Time{}},
	{Title: "有隣堂しか知らない世界", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCmKlo3BXt60nzgk2r_JgvwQ", Fetched: time.Time{}},
	{Title: "ぽんこつ鳩子", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCOs-jFefgxrwMZhK7AcjdMg", Fetched: time.Time{}},
	{Title: "ライアン鈴木", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCYxPWII5Kj5bCRuQ0whNN-Q", Fetched: time.Time{}},
	{Title: "mikimiki", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UChxtIA33ty53Hh4MmkXNASg", Fetched: time.Time{}},
	{Title: "LayerX 公式", URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCKu0R6kOcqd62QstxKxoNhA", Fetched: time.Time{}},

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
