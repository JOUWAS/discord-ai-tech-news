package usecase

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func HandleDiscordMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message from %s: %s", m.Author.Username, m.Content)

	if m.Content == "hello" {
		_, err := s.ChannelMessageSend(m.ChannelID, "sini kau hanip!!!!")
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}
