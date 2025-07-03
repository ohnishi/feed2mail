package common

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ohnishi/feed/config"
	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func MailNotify(msg, attachmentPath string) error {
	appConf, err := config.NewDefaultApp()
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	// メッセージの構築
	message := mail.NewV3Mail()
	// 送信元を設定
	from := mail.NewEmail("", appConf.SendGrid.From)
	message.SetFrom(from)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()
	to := mail.NewEmail("", appConf.SendGrid.To)
	p.AddTos(to)
	p.SetSubstitution("%fullname%", "田中 太郎")
	p.SetSubstitution("%familyname%", "田中")
	p.SetSubstitution("%place%", "中野")
	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = "feed更新通知"
	// テキストパートを設定
	c := mail.NewContent("text/plain", msg)
	message.AddContent(c)

	// カテゴリ情報を付加
	message.AddCategories("notify")
	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// 画像ファイルを添付
	if attachmentPath != "" {
		a := mail.NewAttachment()
		file, err := os.OpenFile(attachmentPath, os.O_RDONLY, 0600)
		if err != nil {
			return errors.Wrap(err, "failed to read torrent.json")
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return errors.Wrap(err, "failed to read all torrent.json")
		}
		data_enc := base64.StdEncoding.EncodeToString(data)
		a.SetContent(data_enc)
		a.SetType("application/json")
		a.SetFilename("torrent.json")
		a.SetDisposition("attachment")
		message.AddAttachment(a)
	}
	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(appConf.SendGrid.APIKey)
	_, err = client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func MailNotifyByResend(msg, attachmentPath string) error {
	appConf, err := config.NewDefaultApp()
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	// apiKey := "re_6vMqHq34_32T3fGeL47rmN43zCLMEZ6MF"

	client := resend.NewClient(appConf.SendGrid.APIKey)

	params := &resend.SendEmailRequest{
		From:    "feed@resend.dev",
		To:      []string{"notify@example.com"},
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

	sent, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	fmt.Println(sent)
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
