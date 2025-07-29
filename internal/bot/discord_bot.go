package bot

import (
	"log"

	"discord-ai-tech-news/internal/handler/discord"

	"github.com/bwmarrin/discordgo"
)

func NewDiscordBot(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	dg.AddHandler(discord.OnMessageCreate)

	if err = dg.Open(); err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}

	return dg
}
