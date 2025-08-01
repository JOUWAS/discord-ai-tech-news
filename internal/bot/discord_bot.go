package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler interface {
	HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate)
}

type DiscordBot struct {
	session *discordgo.Session
}

func NewDiscordBot(token string, handler MessageHandler) *DiscordBot {
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

	return &DiscordBot{
		session: dg,
	}
}

// Close method untuk graceful shutdown
func (bot *DiscordBot) Close() error {
	return bot.session.Close()
}

// SendNewsToChannel mengirim pesan berita ke channel tertentu
func (bot *DiscordBot) SendNewsToChannel(channelName string, message string) error {
	// Cari channel berdasarkan nama
	for _, guild := range bot.session.State.Guilds {
		for _, channel := range guild.Channels {
			if channel.Name == channelName && channel.Type == discordgo.ChannelTypeGuildText {
				_, err := bot.session.ChannelMessageSend(channel.ID, message)
				if err != nil {
					return fmt.Errorf("failed to send message to channel %s: %v", channelName, err)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("channel %s not found", channelName)
}
