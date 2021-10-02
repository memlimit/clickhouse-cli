package config

import (
	"flag"
	"github.com/spf13/viper"
)

type Protocol string

const (
	HTTP Protocol = "http"
)

type Config struct {
	Auth        AuthData `mapstructure:"auth"`
	UseProtocol Protocol `mapstructure:"protocol"`
	HTTP        Http     `mapstructure:"http"`
	CLI         CLI      `mapstructure:"cli"`
}

type CLI struct {
	Multiline   bool   `mapstructure:"multiline"`
	HistoryPath string `mapstructure:"historyPath"`
}

type AuthData struct {
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Http struct {
	URL      string `mapstructure:"url"`
	Compress string `mapstructure:"compress"`
}

func New() (*Config, error) {
	var c Config
	var path string

	flag.StringVar(&path, "config", "$HOME/clickhouse-cli/cli/config/", "set path to config")
	flag.Parse()

	viper.AddConfigPath(path)
	viper.SetConfigName(".clickhouse-cli-config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
