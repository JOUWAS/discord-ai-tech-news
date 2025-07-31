package discord

import (
	"context"
	"log"

	"discord-ai-tech-news/internal/usecase"

	"github.com/bwmarrin/discordgo"
)

// Legacy function - untuk backward compatibility
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

type MessageHandler struct {
	usecase *usecase.MessageUsecase
}

func NewMessageHandler(usecase *usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		usecase: usecase,
	}
}

func (h *MessageHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages from other bots
	if m.Author.Bot {
		return
	}

	// Check if message is in the correct channel
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Failed to get channel: %v", err)
		return
	}

	// Only respond in specific channel
	if channel.Name != "🔥┃ai-tech-news" {
		return
	}

	// Process the message
	ctx := context.Background()
	response, err := h.usecase.ProcessMessage(ctx, m.Content)

	if err != nil {
		log.Printf("Error processing message from %s: %v", m.Author.Username, err)
		response = "❌ **Terjadi kesalahan sistem**\n\n🔄 Silakan coba lagi dalam beberapa saat."
	}

	// Log the interaction
	log.Printf("User %s (%s) sent: %s", m.Author.Username, m.Author.ID, m.Content)

	// Send response
	_, err = s.ChannelMessageSend(m.ChannelID, response)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
