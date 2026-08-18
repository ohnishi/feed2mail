package models

import (
	"log"
	"os"
	"time"

	"github.com/resend/resend-go/v2"
)

const mailRetryInterval = 10 * time.Second

func mailNotifyByResend(msg string, retry int) error {
	apiKey := os.Getenv("MAIL_APIKEY")
	sendFrom := "feed@resend.dev"
	sendTo := "notify@example.com"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    sendFrom,
		To:      []string{sendTo},
		Subject: "feed更新通知",
		Text:    msg,
	}

	// retry が 0 以下でも必ず 1 回は送信を試みる。
	if retry < 1 {
		retry = 1
	}

	var err error
	for cnt := 1; cnt <= retry; cnt++ {
		if _, err = client.Emails.Send(params); err == nil {
			return nil
		}
		if cnt < retry {
			time.Sleep(mailRetryInterval * time.Duration(cnt))
		}
	}
	log.Printf("mail send failed after %d attempts: %v", retry, err)
	return err
}
