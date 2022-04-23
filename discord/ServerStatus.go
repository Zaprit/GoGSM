package discord

import (
	"GoGSM/config"
	"GoGSM/gamequery"
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

const zws = "\u200B"

var spacer = discordgo.MessageEmbedField{
	Name:   zws,
	Value:  zws,
	Inline: true,
}

var goGSMEmbedFooter = discordgo.MessageEmbedFooter{
	Text: "GoGSM 1.0",
}

// RefreshServerStatus refreshes the status embeds for all servers
// TODO: This is too monolithic. It should be broken up into smaller functions
func RefreshServerStatus(s *discordgo.Session) {

	// Refreshes the Game names if they have expired
	config.GetGameNames()

	for name, server := range config.ReadConfig().Servers {

		var statusString = ":yellow_circle: Unknown"
		var messageToEdit string
		var embedFields []*discordgo.MessageEmbedField

		var isOnline bool

		resp, err := gamequery.ServerQuery(&server)

		if errors.Is(gamequery.ErrServerNotFound, err) {
			statusString = ":red_circle: Server not found"
			isOnline = false
		} else if err != nil {
			statusString = ":red_circle: Error querying server"
			isOnline = false
		} else {
			statusString = ":green_circle: Online"
			isOnline = true
		}
		c, er2 := s.ChannelMessages(server.ChannelID, 20, "", "", "")
		if er2 == nil {
			log.Println("Failed to get messages: ", er2)
		}
		for _, i2 := range c {
			if i2.Author.ID != s.State.User.ID {
				continue
			}

			if message, ok := config.GlobalCache.Get(name); ok {
				messageToEdit = message.(string)
			} else {
				err := s.ChannelMessageDelete(server.ChannelID, i2.ID)
				if err != nil {
					log.Println("Error deleting message:", err)
				}
			}
		}

		if isOnline {
			if resp.Name == "" {
				resp.Name = "???"
			}
			maxPlayers := "?"
			if resp.MaxPlayers == 0 {
				maxPlayers = "?"
			} else {
				maxPlayers = strconv.Itoa(resp.MaxPlayers)
			}
			game := server.Game
			if prettyGame, ok := config.PrettyGameNames[server.Game]; ok {
				game = prettyGame
			}

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Name",
				Value:  resp.Name,
				Inline: true,
			})
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Game",
				Value:  game,
				Inline: true,
			})

			embedFields = append(embedFields, &spacer)

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Players",
				Value:  strconv.Itoa(len(resp.Players)) + "/" + maxPlayers,
				Inline: true,
			})
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Map",
				Value:  "`" + resp.Map + "`",
				Inline: true,
			})

			embedFields = append(embedFields, &spacer)

			hostname := server.Hostname
			port := server.Port
			if server.PublicHostname != "" {
				hostname = server.PublicHostname
			}
			if server.PublicPort != 0 {
				port = server.PublicPort
			}

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Country",
				Value:  GetFlagFromCountryCode(server.Country),
				Inline: true,
			})

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Address",
				Value:  "`" + hostname + ":" + strconv.Itoa(port) + "`",
				Inline: true,
			})

			embedFields = append(embedFields, &spacer)

			if server.SteamID != "" {
				if server.DirectJoin {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Join",
						Value:  "steam://connect/" + hostname + ":" + strconv.Itoa(port),
						Inline: true,
					})
				} else {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Join",
						Value:  "steam://rungameid/" + server.SteamID,
						Inline: true,
					})
				}
			}
		}
		var embed = &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       server.Name,
			Description: statusString,
			Timestamp:   time.Now().Format(time.RFC3339),
			Color:       server.Colour,
			Footer:      &goGSMEmbedFooter,
			Thumbnail:   nil,
			Video:       nil,
			Provider:    nil,
			Fields:      embedFields,
		}
		if messageToEdit != "" {
			_, er2 := s.ChannelMessageEditEmbed(server.ChannelID, messageToEdit, embed)
			if er2 != nil {
				log.Println("Error sending embed: ", er2.Error())
			}
		} else {
			m, er2 := s.ChannelMessageSendEmbed(server.ChannelID, embed)
			if er2 != nil {
				log.Println("Error sending embed: ", er2.Error())
			} else {
				config.GlobalCache.Set(name, m.ID, time.Hour*24)
			}
		}

	}

}
