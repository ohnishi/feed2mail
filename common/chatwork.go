package common

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ohnishi/feed/config"
	"github.com/pkg/errors"
)

func postMessageToChatwork(apiToken, roomID, message string) error {
	url := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%s/messages", roomID)
	data := "body=" + message

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("X-ChatWorkToken", apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message posted successfully")
	} else {
		return fmt.Errorf("failed to post message, status code: %d", resp.StatusCode)
	}
	return nil
}

func uploadFileToChatwork(apiToken, roomID, filePath, message string) error {
	url := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%s/files", roomID)

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// multipart/form-data形式でリクエストボディを作成
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// ファイルパートを追加
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// メッセージパートを追加
	err = writer.WriteField("message", message)
	if err != nil {
		return fmt.Errorf("failed to add message field: %v", err)
	}

	// Writerのクローズ
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("X-ChatWorkToken", apiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// リクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスの確認
	if resp.StatusCode == http.StatusOK {
		fmt.Println("File uploaded successfully")
	} else {
		return fmt.Errorf("failed to upload file, status code: %d", resp.StatusCode)
	}
	return nil
}

func ChatworkNotify(message, filePath string) error {
	appConf, err := config.NewDefaultApp()
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	fmt.Println("Uploading file...", appConf.Chatwork.APIToken)

	apiToken := appConf.Chatwork.APIToken
	roomID := "378090501" // 投稿するルームのID

	if filePath == "" {
		// メッセージ投稿
		err = postMessageToChatwork(apiToken, roomID, message)
		if err != nil {
			return errors.Wrap(err, "failed to post message")
		}
	} else {
		// ファイルのアップロード
		err = uploadFileToChatwork(apiToken, roomID, filePath, message)
		if err != nil {
			return errors.Wrap(err, "failed to upload file")
		}
	}

	return nil
}
