package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

// Protocol custom type of int
type Protocol int

const (
	Http Protocol = 0 //nolint:revive
	Grpc Protocol = 1 //nolint:revive
)

// Config - structure of yaml config file
type Config struct {
	Auth     AuthData `mapstructure:"auth"`
	CLI      CLI      `mapstructure:"cli"`
	Protocol Protocol `mapstructure:"protocol"`
	Compress string   `mapstructure:"compress"`
	Address  string   `mapstructure:"address"`
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

// New creates config object with default values or data with file
func New() (*Config, error) {
	var c Config
	var path string

	flag.StringVar(&path, "config", "$HOME/", "set path to config")
	flag.StringVar(&c.Auth.UserName, "username", "default", "set user name")
	flag.StringVar(&c.Auth.Password, "password", "", "set user password")

	flag.StringVar(&c.Address, "address", "127.0.0.1:8123/", "set host:port")
	flag.StringVar(&c.Compress, "compress", "", "set compress method")

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
