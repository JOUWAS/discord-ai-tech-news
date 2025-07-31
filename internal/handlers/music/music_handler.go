package music

import (
	"context"
	"log"
	"strings"

	"discord-ai-tech-news/internal/usecases/music"

	"github.com/bwmarrin/discordgo"
)

type MusicHandler struct {
	usecase *music.MusicUsecase
}

func NewMusicHandler(usecase *music.MusicUsecase) *MusicHandler {
	return &MusicHandler{
		usecase: usecase,
	}
}

func (h *MusicHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages from other bots
	if m.Author.Bot {
		return
	}

	// Only respond to messages starting with ! for music bot
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	// Check if message is in music channels
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Failed to get channel: %v", err)
		return
	}

	// Only respond in music-related channels
	if !strings.Contains(strings.ToLower(channel.Name), "music") &&
		!strings.Contains(strings.ToLower(channel.Name), "bot") &&
		channel.Name != "🕹️┃dev-talk" {
		return
	}

	// Process the message
	ctx := context.Background()
	response, err := h.usecase.ProcessMessage(ctx, m.Content)

	if err != nil {
		log.Printf("Error processing music message from %s: %v", m.Author.Username, err)
		response = "❌ **Music bot error**\n\n🔄 Please try again."
	}

	// Log the interaction
	log.Printf("🎵 Music Bot - User %s (%s) sent: %s", m.Author.Username, m.Author.ID, m.Content)

	// Send response (only if there's a response)
	if response != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Printf("Failed to send music message: %v", err)
		}
	}
}
