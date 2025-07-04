package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/resend/resend-go/v2"
)

func mailNotifyByResend(msg, attachmentPath string) error {
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

	if attachmentPath != "" {
		attachment, err := createAttachment(attachmentPath)
		if err != nil {
			return err
		}
		params.Attachments = append(params.Attachments, attachment)
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("mail send failed: %v", err)
		return err
	}
	return nil
}

// ファイルを読み込んで添付ファイルを作成する関数
func createAttachment(filePath string) (*resend.Attachment, error) {
	// ファイルの内容を読み込む
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	// ファイル名だけを取得
	fileName := filepath.Base(filePath)

	// 添付ファイルを作成
	attachment := &resend.Attachment{
		Filename: fileName,
		Content:  fileContent,
	}

	return attachment, nil
}
