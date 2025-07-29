package discord

import (
	"log"

	"discord-ai-tech-news/internal/usecase"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Failed to get channel: %v", err)
		return
	}

	if channel.Name != "🔥┃ai-tech-news" {
		return
	}

	usecase.HandleDiscordMessage(s, m)
}
