package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
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

	flag.StringVar(&path, "c", "$HOME/", "set path to config")
	flag.StringVar(&c.Auth.UserName, "u", "default", "set user name")
	flag.StringVar(&c.Auth.Password, "p", "", "set user password")
	flag.StringVar(&c.HTTP.URL, "h", "http://127.0.0.1:8123/", "set http host")
	flag.StringVar(&c.HTTP.Compress, "cp", "gzip", "set compress method")
	flag.Parse()

	viper.AddConfigPath(path)
	viper.SetConfigName(".clickhouse-cli-config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config not found at: %s. Set default values.", path)
	}
	_ = viper.Unmarshal(&c)

	return &c, nil
}
