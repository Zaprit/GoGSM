package main

import (
	"GoGSM/config"
	"GoGSM/discord"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "embed"
)

//go:embed config/default.toml
var defaultConfig string

func main() {

	config.DefaultConfig = defaultConfig

	config.GlobalCache = cache.New(5*time.Minute, 10*time.Minute)

	conf := config.ReadConfig()

	// Download names from gamedig repo
	// TODO: this should be cached
	go config.GetGameNames()

	// Start discord bot
	s, err := discord.CreateBot(conf.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating bot: %s", err)
	}

	// Sets the timer to do status updates
	ticker := time.NewTicker(time.Second * time.Duration(conf.Bot.UpdateTime))

	// Does the updates every tick
	go func() {
		for range ticker.C {
			discord.RefreshServerStatus(s)
		}
	}()

	discord.RefreshServerStatus(s)

	// Clean exit on ctrl+c
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// TODO: write cache to disk here
	fmt.Println("\nShutting down...")

}
