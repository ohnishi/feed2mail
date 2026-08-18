// Package feed は RSS/Atom フィードの取得を担当します。
// gofeed への依存をこのパッケージに閉じ込め、上位層には Item/Feed だけを渡します。
package feed

import (
	"context"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Item はフィード中の 1 記事です。
// Published がゼロ値の場合、公開時刻が判定できなかったことを表します。
type Item struct {
	Title     string
	Link      string
	Published time.Time
}

// HasPublished は公開時刻が判定できたかどうかを返します。
func (i Item) HasPublished() bool { return !i.Published.IsZero() }

// Feed は取得済みのフィード 1 件分です。
type Feed struct {
	Items []Item
}

// Fetcher はフィードを取得します。テストでは差し替え可能です。
type Fetcher interface {
	Fetch(ctx context.Context, url string) (*Feed, error)
}

const (
	// ChromeUserAgent は素の Go クライアントを弾くサイト向けの偽装 UA です。
	ChromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

	// DefaultTimeout は接続確立から応答受信までの合計タイムアウトです。
	DefaultTimeout = 30 * time.Second

	// DefaultRetryInterval はリトライ間隔の基準値です。実際の待機は回数に比例して伸びます。
	DefaultRetryInterval = 4 * time.Second
)

type userAgentTransport struct {
	rt http.RoundTripper
}

func (u *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	// Chrome ブラウザからのアクセスであるように偽装するためのヘッダー群
	req.Header.Set("User-Agent", ChromeUserAgent)
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

// GofeedFetcher は gofeed を使う Fetcher の実装です。
type GofeedFetcher struct {
	parser        *gofeed.Parser
	retry         int
	retryInterval time.Duration
	// sleep はリトライ間の待機です。テストから差し替えられるようにしています。
	sleep func(time.Duration)
}

// Option は GofeedFetcher の設定を変更します。
type Option func(*GofeedFetcher)

// WithRetryInterval はリトライ間隔の基準値を変更します。
func WithRetryInterval(d time.Duration) Option {
	return func(f *GofeedFetcher) { f.retryInterval = d }
}

// WithSleep は待機処理を差し替えます。テスト用です。
func WithSleep(sleep func(time.Duration)) Option {
	return func(f *GofeedFetcher) { f.sleep = sleep }
}

// WithTimeout は HTTP クライアントのタイムアウトを変更します。
func WithTimeout(d time.Duration) Option {
	return func(f *GofeedFetcher) { f.parser.Client.Timeout = d }
}

// NewGofeedFetcher は Chrome を騙る HTTP クライアントを備えた Fetcher を返します。
// retry が 1 未満の場合は 1 に丸められます。
func NewGofeedFetcher(retry int, opts ...Option) *GofeedFetcher {
	if retry < 1 {
		retry = 1
	}

	httpClient := &http.Client{
		Timeout:   DefaultTimeout,
		Transport: &userAgentTransport{rt: http.DefaultTransport},
	}

	parser := gofeed.NewParser()
	parser.UserAgent = ChromeUserAgent
	parser.Client = httpClient

	f := &GofeedFetcher{
		parser:        parser,
		retry:         retry,
		retryInterval: DefaultRetryInterval,
		sleep:         time.Sleep,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Fetch はフィードを取得します。失敗した場合は待機時間を増やしながら retry 回まで再試行します。
func (f *GofeedFetcher) Fetch(ctx context.Context, url string) (*Feed, error) {
	var err error
	for cnt := 1; cnt <= f.retry; cnt++ {
		var parsed *gofeed.Feed
		parsed, err = f.parser.ParseURLWithContext(url, ctx)
		if err == nil {
			return convert(parsed), nil
		}
		if ctx.Err() != nil {
			return nil, errors.Wrapf(ctx.Err(), "failed to fetch feed: %s", url)
		}
		if cnt < f.retry {
			f.sleep(f.retryInterval * time.Duration(cnt)) // リトライごとに待機時間を増加
		}
	}
	return nil, errors.Wrapf(err, "failed to fetch feed: %s", url)
}

func convert(parsed *gofeed.Feed) *Feed {
	out := &Feed{Items: make([]Item, 0, len(parsed.Items))}
	for _, item := range parsed.Items {
		out.Items = append(out.Items, Item{
			Title:     item.Title,
			Link:      item.Link,
			Published: publishedAt(item),
		})
	}
	return out
}

// publishedAt は記事の公開時刻を返します。
// gofeed がパース済みの値を使い、published を持たないフィードでは updated で代用します。
// どちらも無い場合はゼロ値を返します。
func publishedAt(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}
	return time.Time{}
}
