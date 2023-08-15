package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fredmayer/mail-parser-rest/internal/configs"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.toml", "path to config file")
}

func main() {
	log.Println("Starting...")
	flag.Parse()

	config := configs.NewConfig(configPath)
	fmt.Println(config)
}
