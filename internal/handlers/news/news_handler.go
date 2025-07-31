package news

import (
	"context"
	"log"

	"discord-ai-tech-news/internal/usecases/news"

	"github.com/bwmarrin/discordgo"
)

type NewsHandler struct {
	usecase *news.NewsUsecase
}

func NewNewsHandler(usecase *news.NewsUsecase) *NewsHandler {
	return &NewsHandler{
		usecase: usecase,
	}
}

func (h *NewsHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	// Only respond in specific channels for news bot
	if channel.Name != "🔥┃ai-tech-news" && channel.Name != "🕹️┃dev-talk" {
		return
	}

	// Process the message
	ctx := context.Background()
	response, err := h.usecase.ProcessMessage(ctx, m.Content)

	if err != nil {
		log.Printf("Error processing news message from %s: %v", m.Author.Username, err)
		response = "❌ **Terjadi kesalahan sistem**\n\n🔄 Silakan coba lagi dalam beberapa saat."
	}

	// Log the interaction
	log.Printf("📰 News Bot - User %s (%s) sent: %s", m.Author.Username, m.Author.ID, m.Content)

	// Send response
	if response != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Printf("Failed to send news message: %v", err)
		}
	}
}
