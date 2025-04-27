package main

import (
	"path/filepath"
	"time"
)

var subscriptions = []subscription{
	{Title: "はてなブックマーク", URL: "https://b.hatena.ne.jp/hotentry/it.rss", Fetched: time.Time{}},
	{Title: "Zennのトレンド", URL: "https://zenn.dev/feed", Fetched: time.Time{}},
	{Title: "Publickey", URL: "https://www.publickey1.jp/atom.xml", Fetched: time.Time{}},
	{Title: "ITmedia", URL: "https://rss.itmedia.co.jp/rss/2.0/topstory.xml", Fetched: time.Time{}},
	{Title: "IT - Yahoo!ニュース", URL: "https://news.yahoo.co.jp/rss/categories/it.xml", Fetched: time.Time{}},
	{Title: "NHKニュース", URL: "https://www3.nhk.or.jp/rss/news/cat0.xml", Fetched: time.Time{}},
	{Title: "まめきちまめこニートの日常", URL: "https://mamekichimameko.blog.jp/index.rdf", Fetched: time.Time{}},
	{Title: "くるねこ大和", URL: "https://blog.goo.ne.jp/kuru0214/rss2.xml", Fetched: time.Time{}},
	{Title: "梅屋敷商店街のランダム・ウォーカー", URL: "https://randomwalker.blog.fc2.com/?xml", Fetched: time.Time{}},
	{Title: "メルカリエンジニアリングブログ", URL: "https://engineering.mercari.com/blog/feed.xml", Fetched: time.Time{}},
	{Title: "しょぷーブログ", URL: "https://new-shopuublog.com/feed/", Fetched: time.Time{}},
	{Title: "チェンソーマン", URL: "https://shonenjumpplus.com/rss/series/3270375685341574005", Fetched: time.Time{}},
	{Title: "ワンパンマン", URL: "https://tonarinoyj.jp/rss/series/13932016480028984490", Fetched: time.Time{}},
	{Title: "タワーダンジョン", URL: "https://comic-days.com/rss/series/14079602755256855913", Fetched: time.Time{}},
	{Title: "ダンダダン", URL: "https://shonenjumpplus.com/rss/series/3269632237310729745", Fetched: time.Time{}},
	{Title: "1688Japan", URL: "https://note.com/1688_japan_blog/rss", Fetched: time.Time{}},
	{Title: "Gemini - Google の AI", URL: "https://note.com/google_gemini/rss", Fetched: time.Time{}},
	{Title: "久松剛", URL: "https://note.com/makaibito/rss", Fetched: time.Time{}},
	{Title: "セールモンスター", URL: "https://note.com/_salemonster/rss", Fetched: time.Time{}},
	{Title: "Amazon物販戦略家！わさびん土屋", URL: "https://note.com/wasabinbin/rss", Fetched: time.Time{}},
	{Title: "shi3z", URL: "https://note.com/shi3zblog/rss", Fetched: time.Time{}},
	{Title: "入江 慎吾🚀 ソロプレナー", URL: "https://note.com/iritec/rss", Fetched: time.Time{}},
	{Title: "tamayan", URL: "https://note.com/tamayan888/rss", Fetched: time.Time{}},
	{Title: "四季報写経ウーマン", URL: "https://note.com/shikiho_shakyo/rss", Fetched: time.Time{}},
	{Title: "AI FREAK - 最新のAIツールをご紹介", URL: "https://note.com/ai_freak/rss", Fetched: time.Time{}},
	{Title: "1688Japan", URL: "https://note.com/1688_japan_blog/rss", Fetched: time.Time{}},
	{Title: "キムラ ヨシト", URL: "https://note.com/k1mu/rss", Fetched: time.Time{}},
	{Title: "ニケちゃん", URL: "https://note.com/nike_cha_n/rss", Fetched: time.Time{}},
	{Title: "erukiti", URL: "https://note.com/erukiti/rss", Fetched: time.Time{}},
	{Title: "shu", URL: "https://note.com/shu127/rss", Fetched: time.Time{}},
	{Title: "現海秀哲物販Note", URL: "https://note.com/genkai_hideaki/rss", Fetched: time.Time{}},
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
			s.Fetched = time.Now().AddDate(0, 0, -1)
			newSubscriptions = append(newSubscriptions, s)
		}
	}

	return writeSubscription(filePath, newSubscriptions)
}
