package config

type SendGridApp struct {
	APIKey string `mapstructure:"apikey"`
	To     string `mapstructure:"to"`
	From   string `mapstructure:"from"`
}
