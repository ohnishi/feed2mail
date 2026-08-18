// Package store は購読状態の永続化を担当します。
// 差分が読みやすいよう、いずれも 1 レコード 1 行の JSONL 形式で保存します。
package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// Subscription は購読するフィード 1 件と、その最終取得時刻です。
type Subscription struct {
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Fetched time.Time `json:"fetched"`
}

// LoadSubscriptions は JSONL 形式の購読状態を読み込みます。
// ファイルが存在しない場合は空のリストを返します。
func LoadSubscriptions(path string) ([]Subscription, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var subscriptions []Subscription
	d := json.NewDecoder(f)
	for d.More() {
		var s Subscription
		if err := d.Decode(&s); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal subscription: %s", path)
		}
		subscriptions = append(subscriptions, s)
	}
	return subscriptions, nil
}

// SaveSubscriptions は購読状態を JSONL 形式で書き出します。
func SaveSubscriptions(path string, subscriptions []Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}
	return writeJSONL(path, func(enc *json.Encoder) error {
		for _, s := range subscriptions {
			if err := enc.Encode(s); err != nil {
				return errors.Wrapf(err, "failed to write subscription: %s", s.URL)
			}
		}
		return nil
	})
}

// writeJSONL は path を作り直して JSONL を書き込み、fsync まで行います。
func writeJSONL(path string, write func(*json.Encoder) error) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.Wrapf(err, "failed to create output directory: %s", dir)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	if err := write(json.NewEncoder(f)); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return errors.Wrapf(err, "failed to sync file: %s", path)
	}
	return nil
}
