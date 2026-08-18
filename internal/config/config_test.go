package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const minimalYAML = `
subscriptions:
  - title: "Zenn"
    url: "https://zenn.dev/feed"
`

// writeConfig は一時ディレクトリに設定ファイルを書き出してパスを返します。
func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

// setMailEnv は検証を通過する最低限の環境変数を設定します。
func setMailEnv(t *testing.T, to string) {
	t.Helper()
	t.Setenv(EnvAPIKey, "test-key")
	t.Setenv(EnvMailTo, to)
}

func TestLoadAppliesDefaults(t *testing.T) {
	setMailEnv(t, "someone@example.com")

	cfg, err := Load(writeConfig(t, minimalYAML))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	def := Default()
	if cfg.StatePath != def.StatePath {
		t.Errorf("StatePath = %q, want %q", cfg.StatePath, def.StatePath)
	}
	if cfg.SeenRetention != def.SeenRetention {
		t.Errorf("SeenRetention = %v, want %v", cfg.SeenRetention, def.SeenRetention)
	}
	if cfg.MaxRetry != def.MaxRetry {
		t.Errorf("MaxRetry = %d, want %d", cfg.MaxRetry, def.MaxRetry)
	}
	if cfg.Mail.Subject != def.Mail.Subject {
		t.Errorf("Mail.Subject = %q, want %q", cfg.Mail.Subject, def.Mail.Subject)
	}
}

func TestLoadParsesDurations(t *testing.T) {
	setMailEnv(t, "someone@example.com")

	cfg, err := Load(writeConfig(t, minimalYAML+"\nseen_retention: 168h\npolite_delay: 500ms\n"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.SeenRetention != 168*time.Hour {
		t.Errorf("SeenRetention = %v, want 168h", cfg.SeenRetention)
	}
	if cfg.PoliteDelay != 500*time.Millisecond {
		t.Errorf("PoliteDelay = %v, want 500ms", cfg.PoliteDelay)
	}
}

// 宛先はリポジトリに残さないため、環境変数だけから読む。
func TestLoadReadsRecipientsFromEnv(t *testing.T) {
	setMailEnv(t, " a@example.com , b@example.com ,")
	t.Setenv(EnvMailFrom, "sender@example.com")

	cfg, err := Load(writeConfig(t, minimalYAML))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := []string{"a@example.com", "b@example.com"}
	if len(cfg.Mail.To) != len(want) {
		t.Fatalf("Mail.To = %v, want %v", cfg.Mail.To, want)
	}
	for i := range want {
		if cfg.Mail.To[i] != want[i] {
			t.Errorf("Mail.To[%d] = %q, want %q", i, cfg.Mail.To[i], want[i])
		}
	}
	if cfg.Mail.From != "sender@example.com" {
		t.Errorf("Mail.From = %q, want sender@example.com", cfg.Mail.From)
	}
}

func TestLoadValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		apiKey string
		mailTo string
	}{
		{name: "購読リストが空", body: "subscriptions: []\n", apiKey: "k", mailTo: "a@example.com"},
		{name: "url が無い", body: "subscriptions:\n  - title: \"Zenn\"\n", apiKey: "k", mailTo: "a@example.com"},
		{name: "APIキー未設定", body: minimalYAML, apiKey: "", mailTo: "a@example.com"},
		{name: "宛先未設定", body: minimalYAML, apiKey: "k", mailTo: ""},
		{name: "保持期間が非正", body: minimalYAML + "\nseen_retention: 0s\n", apiKey: "k", mailTo: "a@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvAPIKey, tt.apiKey)
			t.Setenv(EnvMailTo, tt.mailTo)
			if _, err := Load(writeConfig(t, tt.body)); err == nil {
				t.Error("Load() error = nil, want error")
			}
		})
	}
}

func TestLoadMissingFile(t *testing.T) {
	setMailEnv(t, "a@example.com")
	if _, err := Load(filepath.Join(t.TempDir(), "absent.yaml")); err == nil {
		t.Error("Load() error = nil, want error")
	}
}

// リポジトリ同梱の config.yaml が実際に読めることを確認する。
func TestLoadRepositoryConfig(t *testing.T) {
	setMailEnv(t, "a@example.com")

	cfg, err := Load(filepath.Join("..", "..", DefaultPath))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Subscriptions) == 0 {
		t.Error("Subscriptions is empty")
	}
	if len(cfg.ExcludePrefixes) == 0 {
		t.Error("ExcludePrefixes is empty")
	}
}
