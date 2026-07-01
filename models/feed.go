// Package models は、このプログラムのエントリポイントを提供します。
package models

import (
	"encoding/json"
	"fmt"
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

type userAgentTransport struct {
	rt http.RoundTripper
}

func (u *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	// Chromeブラウザからのアクセスであるように偽装するためのヘッダー群
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
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

	// 🌟 修正点 1: タイムアウト設定を持つカスタムクライアントを作成 🌟
	// 接続確立から応答受信までの合計タイムアウトを30秒に設定しつつ、Chromeブラウザからのアクセスに偽装する
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &userAgentTransport{
			rt: http.DefaultTransport,
		},
	}

	fp := gofeed.NewParser()
	fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	fp.Client = httpClient // カスタムクライアントをParserに設定

	var bodys []string
	var itemLinkMap = make(map[string]struct{})

	for i, subscription := range subscriptions {
		var latest time.Time
		var feed *gofeed.Feed
		for cnt := 1; cnt <= retry; cnt++ {
			feed, err = fp.ParseURL(subscription.URL)
			// log.Printf("parse rss url: %d, %s", cnt, subscription.URL)
			if err == nil {
				time.Sleep(3 * time.Second * time.Duration(cnt)) // リトライごとに待機時間を増加
				break
			} else if cnt == retry {
				log.Printf("HTTP request failed: %v", err)
				return err
			}
			time.Sleep(10 * time.Second)
		}

		siteBodys := []string{subscription.Title}
		for _, item := range feed.Items {
			if item.Published == "" {
				continue
			}

			published, err := parseLocal(item.Published)
			if err != nil {
				return err
			}

			if latest.Before(published) {
				latest = published
			}

			_, ok := itemLinkMap[item.Link]
			if subscription.Fetched.Before(published) && !ok {
				// item.Link が "https://anond.hatelabo.jp" から始まる文字列ならスキップ
				if strings.HasPrefix(item.Link, "https://anond.hatelabo.jp") {
					continue
				}

				// item.Link が "https://www.youtube.com/shorts/" から始まる文字列ならスキップ
				if strings.HasPrefix(item.Link, "https://www.youtube.com/shorts/") {
					continue
				}

				siteBodys = append(siteBodys, "  - "+item.Title)
				siteBodys = append(siteBodys, "    - "+item.Link)
				itemLinkMap[item.Link] = struct{}{}
			}
		}
		if len(siteBodys) > 1 {
			if len(bodys) > 0 {
				bodys = append(bodys, "")
			}
			bodys = append(bodys, siteBodys...)
		}

		if subscription.Fetched.Before(latest) {
			subscriptions[i].Fetched = latest
		}
	}

	err = mailNotifyByResend(strings.Join(bodys, "\n"), retry)
	if err != nil {
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

	for _, subscription := range subscriptions {

		err = appendOutFile(f, subscription)
		if err != nil {
			return err
		}
	}
	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

func createOutFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create output directory: %s", dir)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	return f, nil
}

func appendOutFile(f *os.File, v interface{}) error {
	jsonl, err := toJSON(v)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(jsonl)); err != nil {
		return errors.Wrapf(err, "failed to write line: %s", jsonl)
	}
	return nil
}
func toJSON(r interface{}) (string, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return "", errors.Wrapf(err, "could not marshal: %v", r)
	}
	return fmt.Sprintf("%s\n", jsonStr), nil
}

func parseLocal(value string) (t time.Time, err error) {
	for _, layout := range []string{time.RFC1123, time.RFC3339, time.RFC1123Z} {
		t, err = time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "cannot parse as %q", time.Local)
	}
	return t, nil
}

// FileExists はファイルが存在するかどうかを確認します。
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func ReadSubscriptions(filePath string) ([]Subscription, error) {
	if !FileExists(filePath) {
		return nil, nil
	}

	f, err := os.Open(filePath)
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
