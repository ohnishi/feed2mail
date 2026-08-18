// Package config は設定ファイルと環境変数から実行時設定を読み込みます。
//
// 購読リストや除外ルールのように編集頻度の高いものは YAML に置き、
// 秘密情報や宛先のような環境依存の値は環境変数から読みます。
package config

import (
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// 環境変数名。
const (
	EnvAPIKey   = "MAIL_APIKEY"
	EnvMailFrom = "MAIL_FROM"
	EnvMailTo   = "MAIL_TO"
)

// DefaultPath は設定ファイルの既定の場所です。
const DefaultPath = "config.yaml"

// Subscription は購読するフィード 1 件です。
type Subscription struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// Mail は通知メールの設定です。宛先と送信元は環境変数で上書きされます。
type Mail struct {
	Subject string `yaml:"subject"`
	From    string `yaml:"from"`
	To      []string
}

// Config は実行時設定の全体です。
type Config struct {
	// StatePath は購読状態を保存する JSONL のパスです。
	StatePath string `yaml:"state_path"`
	// SeenPath は既読リンクを保存する JSONL のパスです。
	SeenPath string `yaml:"seen_path"`
	// SeenRetention は既読リンクを保持する期間です。
	// 長くするほど取りこぼしに強くなりますが、状態ファイルも比例して大きくなります。
	SeenRetention time.Duration `yaml:"seen_retention"`
	// NewSubLookback は新規購読フィードで初回に遡る期間です。
	NewSubLookback time.Duration `yaml:"new_subscription_lookback"`
	// PoliteDelay はフィード間の待機時間です。
	PoliteDelay time.Duration `yaml:"polite_delay"`
	// MaxRetry はフィード取得とメール送信の最大試行回数です。
	MaxRetry int `yaml:"max_retry"`
	// ExcludePrefixes はこの接頭辞で始まるリンクを通知対象から除外します。
	ExcludePrefixes []string `yaml:"exclude_prefixes"`
	// Subscriptions は購読するフィードの一覧です。
	Subscriptions []Subscription `yaml:"subscriptions"`

	Mail Mail `yaml:"mail"`

	// APIKey は Resend の API キーです。設定ファイルには置かず環境変数から読みます。
	APIKey string `yaml:"-"`
}

// Default は設定ファイルで省略された項目に使う既定値を返します。
func Default() Config {
	return Config{
		StatePath:      "cache/fetched.jsonl",
		SeenPath:       "cache/seen.jsonl",
		SeenRetention:  30 * 24 * time.Hour,
		NewSubLookback: 24 * time.Hour,
		PoliteDelay:    2 * time.Second,
		MaxRetry:       5,
		Mail:           Mail{Subject: "feed更新通知", From: "feed@resend.dev"},
	}
}

// Load は設定ファイルを読み込み、環境変数を重ねて検証済みの設定を返します。
func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrapf(err, "failed to read config: %s", path)
	}

	cfg := Default()
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, errors.Wrapf(err, "failed to parse config: %s", path)
	}

	cfg.applyEnv()
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// applyEnv は環境変数で設定を上書きします。
// 宛先や API キーは環境ごとに変わり、リポジトリに残したくないため設定ファイルには置きません。
func (c *Config) applyEnv() {
	c.APIKey = os.Getenv(EnvAPIKey)
	if from := os.Getenv(EnvMailFrom); from != "" {
		c.Mail.From = from
	}
	if to := os.Getenv(EnvMailTo); to != "" {
		c.Mail.To = splitAndTrim(to)
	}
}

func (c Config) validate() error {
	if len(c.Subscriptions) == 0 {
		return errors.New("config: subscriptions is empty")
	}
	for i, s := range c.Subscriptions {
		if s.Title == "" || s.URL == "" {
			return errors.Errorf("config: subscriptions[%d] requires both title and url", i)
		}
	}
	if c.APIKey == "" {
		return errors.Errorf("config: %s is not set", EnvAPIKey)
	}
	if len(c.Mail.To) == 0 {
		return errors.Errorf("config: %s is not set", EnvMailTo)
	}
	if c.Mail.From == "" {
		return errors.Errorf("config: mail.from is empty (set it in the config file or %s)", EnvMailFrom)
	}
	if c.SeenRetention <= 0 {
		return errors.New("config: seen_retention must be positive")
	}
	return nil
}

// splitAndTrim はカンマ区切りの値を分解し、空要素を落とします。
func splitAndTrim(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
