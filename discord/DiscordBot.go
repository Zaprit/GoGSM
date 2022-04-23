package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "refresh",
		Description: "Refreshes Server List",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"refresh": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Refreshing...",
			},
		})
		if err != nil {
			return
		}
		RefreshServerStatus(s)
		er2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Done!",
			},
		})
		if er2 != nil {
			return
		}
	},
}

func CreateBot(token string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil, err
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	// Open a websocket connection to Discord and begin listening.
	er2 := discord.Open()
	if er2 != nil {
		fmt.Println("error opening connection,", err)
		return nil, er2
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	return discord, nil
}

func GetFlagFromCountryCode(country string) string {
	return ":flag_" + strings.ToLower(country) + ":"
}
