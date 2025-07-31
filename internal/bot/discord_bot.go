package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler interface {
	HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewDiscordBot(token string, handler MessageHandler) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Use injected handler instead of hard-coded one
	dg.AddHandler(handler.HandleMessage)

	if err = dg.Open(); err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}

	log.Println("Discord bot connected successfully!")
	return dg
}
