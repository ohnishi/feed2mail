// Package notify は通知の送信を担当します。
package notify

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
)

// Notifier は本文を送信します。テストでは差し替え可能です。
type Notifier interface {
	Notify(ctx context.Context, subject, body string) error
}

// DefaultRetryInterval はリトライ間隔の基準値です。実際の待機は回数に比例して伸びます。
const DefaultRetryInterval = 10 * time.Second

// Resend は Resend API でメールを送る Notifier の実装です。
type Resend struct {
	client        *resend.Client
	from          string
	to            []string
	retry         int
	retryInterval time.Duration
	sleep         func(time.Duration)
}

// Option は Resend の設定を変更します。
type Option func(*Resend)

// WithRetryInterval はリトライ間隔の基準値を変更します。
func WithRetryInterval(d time.Duration) Option {
	return func(r *Resend) { r.retryInterval = d }
}

// WithSleep は待機処理を差し替えます。テスト用です。
func WithSleep(sleep func(time.Duration)) Option {
	return func(r *Resend) { r.sleep = sleep }
}

// NewResend は Resend を使う Notifier を返します。
// retry が 1 未満の場合は 1 に丸められ、必ず 1 回は送信を試みます。
func NewResend(apiKey, from string, to []string, retry int, opts ...Option) *Resend {
	if retry < 1 {
		retry = 1
	}
	r := &Resend{
		client:        resend.NewClient(apiKey),
		from:          from,
		to:            to,
		retry:         retry,
		retryInterval: DefaultRetryInterval,
		sleep:         time.Sleep,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Notify はメールを送信します。失敗した場合は待機時間を増やしながら retry 回まで再試行します。
func (r *Resend) Notify(ctx context.Context, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    r.from,
		To:      r.to,
		Subject: subject,
		Text:    body,
	}

	var err error
	for cnt := 1; cnt <= r.retry; cnt++ {
		if _, err = r.client.Emails.SendWithContext(ctx, params); err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return errors.Wrap(ctx.Err(), "mail send canceled")
		}
		if cnt < r.retry {
			r.sleep(r.retryInterval * time.Duration(cnt))
		}
	}
	log.Printf("mail send failed after %d attempts: %v", r.retry, err)
	return errors.Wrap(err, "failed to send mail")
}
