package config

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// App は全サービスのappの設定を表す。
type App struct {
	SendGrid SendGridApp `mapstructure:"sendgrid"`
}

// NewApp はconfigPathの設定を読んで新しいAppを返す。
func NewApp(path string) (*App, error) {
	c := App{}
	err := readToml(path, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// AppPath は規定の外部サービスへの接続設定ファイルのパスを返す。
func AppPath() string {
	return configDir("feed.toml")
}

// NewDefaultApp はデフォルトのパスから設定を読んで新しいAppConfigを返す。
func NewDefaultApp() (*App, error) {
	return NewApp(AppPath())
}

// ReadToml はTOMLのファイルを読んでdataに格納する。
func readToml(path string, data interface{}) error {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read TOML file")
	}
	err = v.Unmarshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal TOML data")
	}
	return nil
}

// configDir は設定ディレクトリのパスを返す。
func configDir(subpaths ...string) string {
	return filepath.Join(appRootDir(), "config", filepath.Join(subpaths...))
}

// appRootDir はconfig, dataディレクトリを置くためのルートディレクトリ。
// デフォルトは$GOPATH。
// quaと同じGOPATH下でconfig等を切り替えたい場合に指定すること。
func appRootDir() string {
	return getEnv("APP_ROOT_DIR", build.Default.GOPATH)
}

func getEnv(name string, alt string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return alt
}
