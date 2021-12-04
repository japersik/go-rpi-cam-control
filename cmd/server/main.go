package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/japersik/go-rpi-cam-control/internal/app/server"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configs-path", "configs/serverConf.toml", "path to configs file")

}
func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
