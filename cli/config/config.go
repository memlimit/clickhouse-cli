package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

type Protocol int

const (
	Http Protocol = 0
	Grpc Protocol = 1
)

// Config - structure of yaml config file
type Config struct {
	Auth     AuthData `mapstructure:"auth"`
	HTTP     HTTP     `mapstructure:"http"`
	GRPC     GRPC     `mapstructure:"http"`
	CLI      CLI      `mapstructure:"cli"`
	Protocol Protocol `mapstructure:"protocol"`
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

// GRPC config part with url and compression type
type GRPC struct {
	URL      string `mapstructure:"url"`
	Compress string `mapstructure:"compress"`
}

// New creates config object with default values or data with file
func New() (*Config, error) {
	var c Config
	var path string

	flag.StringVar(&path, "config", "$HOME/", "set path to config")
	flag.StringVar(&c.Auth.UserName, "username", "default", "set user name")
	flag.StringVar(&c.Auth.Password, "password", "", "set user password")

	flag.StringVar(&c.HTTP.URL, "http", "http://127.0.0.1:8123/", "set http host")
	flag.StringVar(&c.HTTP.Compress, "h-compress", "", "set compress method for http")

	flag.StringVar(&c.GRPC.URL, "grpc", "127.0.0.1:9100/", "set grpc host")
	flag.StringVar(&c.GRPC.Compress, "g-compress", "", "set compress method for grpc")

	flag.IntVar((*int)(&c.Protocol), "protocol", 0, "set default protocol. http/grpc")
	flag.Parse()

	viper.AddConfigPath(path)
	viper.SetConfigName(".clickhouse-cli-config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config not found at: %s. Set default values.", path)
	}
	_ = viper.Unmarshal(&c)

	return &c, nil
}
