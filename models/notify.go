package models

import (
	"log"
	"os"
	"time"

	"github.com/resend/resend-go/v2"
)

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

	for cnt := 1; cnt <= retry; cnt++ {
		_, err := client.Emails.Send(params)
		if err == nil {
			break
		} else if cnt == retry {
			log.Printf("mail send failed: %v", err)
			return err
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}
