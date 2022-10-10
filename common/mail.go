package common

import (
	"github.com/ohnishi/feed/config"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func MailNotify(msg string) error {
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

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(appConf.SendGrid.APIKey)
	_, err = client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
