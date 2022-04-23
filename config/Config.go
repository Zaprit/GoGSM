package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
)

var DefaultConfig string

var GlobalCache *cache.Cache

// Config is the main configuration struct
type Config struct {
	Bot     BotConfig
	Servers map[string]Server
}

// BotConfig is configuration for the discord bot
type BotConfig struct {
	Token           string
	UpdateTime      int `toml:"update_interval"`
	Prefix          string
	UseRichPresence bool
}

// Server is configuration for a server, is in an array in Config
type Server struct {
	Name           string `toml:"name"`
	Hostname       string `toml:"hostname"`
	Game           string `toml:"game"`
	Port           int    `toml:"port"`
	Colour         int    `toml:"colour"`
	PublicHostname string `toml:"public_hostname"`
	HideMap        bool   `toml:"hide_map"`
	ChannelID      string `toml:"channel"`
	Country        string `toml:"country"`
	SteamID        string `toml:"steamid"`
	PublicPort     int    `toml:"public_port"`
	DirectJoin     bool   `toml:"direct_join"`
}

// ReadConfig creates a Config struct from config.toml and if config.toml doesn't exist it creates it from default.toml
func ReadConfig() Config {
	var conf Config
	_, err := os.Stat("config.toml")
	if os.IsNotExist(err) {
		file, er2 := os.Create("config.toml")
		if er2 != nil {
			fmt.Println("Error creating config file: config.toml")
			panic(er2)
		}
		_, er3 := file.WriteString(DefaultConfig)
		if er3 != nil {
			panic(er3.Error())
		}
		file.Sync()
		er4 := file.Close()
		if er4 != nil {
			panic(er4.Error())
		}
		log.Println("Created config file: config.toml, please fill in your token and other values")
		log.Println(DefaultConfig)
		os.Exit(0)
	}
	_, er2 := toml.DecodeFile("config.toml", &conf)
	if er2 != nil {
		panic(er2.Error())
	}
	return conf
}
