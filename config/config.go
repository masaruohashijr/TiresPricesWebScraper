package config

import (
	m "4tiresWebScraper/models"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var TiresConfig m.Config

func InitConfig() {
	TiresConfig = ReadConfig()
}

func ReadConfig() m.Config {
	var configfile = "config/config.cfg"
	log.Println(configfile)
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("File configuration "+configfile+" missing: ", configfile)
	}
	var config m.Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
