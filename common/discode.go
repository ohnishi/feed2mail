package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type WebhookMessage struct {
	Content string `json:"content"`
}

func sendDiscordWebhook(webhookURL, message, filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// ファイルが存在しない場合、テキストメッセージのみを送信
		return sendTextOnlyWebhook(webhookURL, message)
	}

	// ファイルが存在する場合、メッセージとファイルを含むマルチパートのボディを作成
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// メッセージ部分の追加
	_ = writer.WriteField("content", message)

	// ファイル部分の追加
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("failed to copy file data: %v", err)
	}

	// マルチパートの書き込みを終了
	writer.Close()

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", webhookURL, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// リクエストの送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-200/204 response: %d", resp.StatusCode)
	}

	return nil
}

func sendTextOnlyWebhook(webhookURL, message string) error {
	// Webhookのメッセージを設定
	webhookMessage := WebhookMessage{Content: message}

	// JSONにエンコード
	jsonMessage, err := json.Marshal(webhookMessage)
	if err != nil {
		return fmt.Errorf("failed to encode message to JSON: %v", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-200/204 response: %d", resp.StatusCode)
	}

	return nil
}

func DiscodeNotify2(message, filePath string) error {
	webhookURL := "https://discord.com/api/webhooks/1304513906919870555/N7Iuc1f3Csalc1ZbOZ8jc-Kbn2fheVTak_X_Ph1_nO2HmarK7xTkCVY2pVEaAazIrvvw" // Webhook URLを指定
	return sendDiscordWebhook(webhookURL, message, filePath)
}
