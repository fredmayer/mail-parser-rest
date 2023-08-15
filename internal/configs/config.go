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
}

func NewConfig(path string) *Config {
	c := &Config{}
	path = "../../" + path
	_, err := toml.DecodeFile(path, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
