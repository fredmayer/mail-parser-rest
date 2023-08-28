package configs

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Imap         string `toml:"imap"`
	ImapUser     string `toml:"imap_user"`
	ImapPassword string `toml:"imap_password"`
	ImapPort     string `toml:"imap_port"`
	ImapSsl      bool   `toml:"imap_ssl"`
	HTTPAddr     string `toml:"http_addr"`

	LogLevel string `toml:"log_level"`
}

var (
	config Config
)

func NewConfig(path string) *Config {
	path = "../../" + path
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func Get() *Config {
	return &config
}
