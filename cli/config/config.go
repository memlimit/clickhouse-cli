package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

// Config - structure of yaml config file
type Config struct {
	Auth AuthData `mapstructure:"auth"`
	HTTP HTTP     `mapstructure:"http"`
	CLI  CLI      `mapstructure:"cli"`
}

// CLI config part with path to history file and multiline input state
type CLI struct {
	Multiline   bool   `mapstructure:"multiline"`
	HistoryPath string `mapstructure:"historyPath"`
}

// AuthData config part with username and password
type AuthData struct {
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// HTTP config part with url and compression type
type HTTP struct {
	URL      string `mapstructure:"url"`
	Compress string `mapstructure:"compress"`
}

// New creates config object with default values or data with file
func New() (*Config, error) {
	var c Config
	var path string

	flag.StringVar(&path, "c", "$HOME/", "set path to config")
	flag.StringVar(&c.Auth.UserName, "u", "default", "set user name")
	flag.StringVar(&c.Auth.Password, "p", "", "set user password")
	flag.StringVar(&c.HTTP.URL, "h", "http://127.0.0.1:8123/", "set http host")
	flag.StringVar(&c.HTTP.Compress, "cp", "", "set compress method")
	flag.Parse()

	viper.AddConfigPath(path)
	viper.SetConfigName(".clickhouse-cli-config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config not found at: %s. Set default values.", path)
	}
	_ = viper.Unmarshal(&c)

	return &c, nil
}
