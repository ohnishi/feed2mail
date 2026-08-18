// Package models は、フィードの取得・重複排除・メール通知・状態の永続化を提供します。
package models

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type Subscription struct {
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Fetched time.Time `json:"fetched"`
}

const (
	chromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	fetchTimeout    = 30 * time.Second
	retryInterval   = 4 * time.Second
	// politeDelay はフィード間の待機時間です。
	// YouTube のフィードは短時間に連続アクセスすると 404/500 を返すことがあるため、間隔を空ける。
	politeDelay = 2 * time.Second
)

type userAgentTransport struct {
	rt http.RoundTripper
}

func (u *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	// Chromeブラウザからのアクセスであるように偽装するためのヘッダー群
	req.Header.Set("User-Agent", chromeUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)

	rt := u.rt
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(req)
}

func FeedToMail(filePath string, retry int) error {
	subscriptions, err := ReadSubscriptions(filePath)
	if err != nil {
		return err
	}

	// 接続確立から応答受信までの合計タイムアウトを設定しつつ、Chromeブラウザからのアクセスに偽装する
	httpClient := &http.Client{
		Timeout: fetchTimeout,
		Transport: &userAgentTransport{
			rt: http.DefaultTransport,
		},
	}

	fp := gofeed.NewParser()
	fp.UserAgent = chromeUserAgent
	fp.Client = httpClient // カスタムクライアントをParserに設定

	var bodies []string
	var failed []string
	var itemLinkMap = make(map[string]struct{})

	for i, subscription := range subscriptions {
		if i > 0 {
			time.Sleep(politeDelay)
		}

		feed, err := fetchFeed(fp, subscription.URL, retry)
		if err != nil {
			// 1つのフィードの失敗で全体を止めず、取得できた分だけ通知する。
			// Fetched も更新しないので、次回の実行で取りこぼしを拾い直せる。
			log.Printf("skip %s: %v", subscription.Title, err)
			failed = append(failed, subscription.Title)
			continue
		}

		now := time.Now()
		var latest time.Time
		siteBodies := []string{subscription.Title}
		for _, item := range feed.Items {
			published, ok := publishedAt(item)
			if !ok {
				continue
			}

			// 未来日付のアイテムで Fetched が先に進むと、以降の新着を恒久的に
			// 取りこぼすため、現在時刻を超える公開時刻は latest に反映しない。
			if published.After(now) {
				continue
			}

			if latest.Before(published) {
				latest = published
			}

			_, sent := itemLinkMap[item.Link]
			if subscription.Fetched.Before(published) && !sent {
				// item.Link が "https://anond.hatelabo.jp" から始まる文字列ならスキップ
				if strings.HasPrefix(item.Link, "https://anond.hatelabo.jp") {
					continue
				}

				// item.Link が "https://www.youtube.com/shorts/" から始まる文字列ならスキップ
				if strings.HasPrefix(item.Link, "https://www.youtube.com/shorts/") {
					continue
				}

				siteBodies = append(siteBodies, "  - "+item.Title)
				siteBodies = append(siteBodies, "    - "+item.Link)
				itemLinkMap[item.Link] = struct{}{}
			}
		}
		if len(siteBodies) > 1 {
			if len(bodies) > 0 {
				bodies = append(bodies, "")
			}
			bodies = append(bodies, siteBodies...)
		}

		if subscription.Fetched.Before(latest) {
			subscriptions[i].Fetched = latest
		}
	}

	if len(failed) > 0 {
		if len(bodies) > 0 {
			bodies = append(bodies, "")
		}
		bodies = append(bodies, "取得に失敗したフィード:")
		for _, title := range failed {
			bodies = append(bodies, "  - "+title)
		}
	}

	// 新着も失敗もない場合は空メールを送らない。
	if len(bodies) == 0 {
		log.Print("no new items; skip notification")
		return WriteSubscription(filePath, subscriptions)
	}

	if err := mailNotifyByResend(strings.Join(bodies, "\n"), retry); err != nil {
		return err
	}

	return WriteSubscription(filePath, subscriptions)
}

func WriteSubscription(path string, subscriptions []Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}
	f, err := createOutFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, subscription := range subscriptions {
		if err := enc.Encode(subscription); err != nil {
			return errors.Wrapf(err, "failed to write subscription: %s", subscription.URL)
		}
	}
	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

func createOutFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create output directory: %s", dir)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	return f, nil
}

// fetchFeed はフィードを取得します。失敗した場合は待機時間を増やしながら retry 回まで再試行します。
func fetchFeed(fp *gofeed.Parser, url string, retry int) (*gofeed.Feed, error) {
	if retry < 1 {
		retry = 1
	}

	var err error
	for cnt := 1; cnt <= retry; cnt++ {
		var feed *gofeed.Feed
		feed, err = fp.ParseURL(url)
		if err == nil {
			return feed, nil
		}
		if cnt < retry {
			time.Sleep(retryInterval * time.Duration(cnt)) // リトライごとに待機時間を増加
		}
	}
	return nil, errors.Wrapf(err, "failed to fetch feed: %s", url)
}

// publishedAt は記事の公開時刻を返します。
// gofeed がパース済みの値を使い、published を持たないフィードでは updated で代用します。
func publishedAt(item *gofeed.Item) (time.Time, bool) {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed, true
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed, true
	}
	return time.Time{}, false
}

// ReadSubscriptions は JSONL 形式の購読リストを読み込みます。
// ファイルが存在しない場合は空のリストを返します。
func ReadSubscriptions(filePath string) ([]Subscription, error) {
	f, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", filePath)
	}
	defer f.Close()

	var subscriptions []Subscription
	d := json.NewDecoder(f)
	for d.More() {
		var s Subscription
		if err := d.Decode(&s); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal subscription")
		}
		subscriptions = append(subscriptions, s)
	}
	return subscriptions, nil
}
