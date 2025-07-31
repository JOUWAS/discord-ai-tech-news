package music

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type MusicUsecase struct {
	// musicService music.MusicService // Will be implemented later
}

func NewMusicUsecase() *MusicUsecase {
	return &MusicUsecase{}
}

func (u *MusicUsecase) ProcessMessage(ctx context.Context, content string) (string, error) {
	content = strings.TrimSpace(content)

	// Only respond to commands starting with !
	if !strings.HasPrefix(content, "!") {
		return "", nil // Don't respond
	}

	// Remove ! prefix
	command := strings.ToLower(strings.TrimPrefix(content, "!"))
	parts := strings.Split(command, " ")

	switch parts[0] {
	case "play", "p":
		if len(parts) < 2 {
			return "❌ Format: `!play <song name or URL>`", nil
		}
		query := strings.Join(parts[1:], " ")
		return u.handlePlay(ctx, query)
	case "stop":
		return u.handleStop(ctx)
	case "pause":
		return u.handlePause(ctx)
	case "resume":
		return u.handleResume(ctx)
	case "queue", "q":
		return u.handleQueue(ctx)
	case "skip":
		return u.handleSkip(ctx)
	case "help":
		return u.getHelpMessage(), nil
	default:
		return "", nil // Don't respond to unknown commands
	}
}

func (u *MusicUsecase) handlePlay(ctx context.Context, query string) (string, error) {
	// TODO: Implement actual music playing
	log.Printf("🎵 DEBUG: User wants to play: %s", query)
	return fmt.Sprintf("🎵 **Mock Response**: Would play \"%s\"\n\n⚠️ Music functionality will be implemented next!", query), nil
}

func (u *MusicUsecase) handleStop(ctx context.Context) (string, error) {
	return "⏹️ **Music stopped** (Mock response)", nil
}

func (u *MusicUsecase) handlePause(ctx context.Context) (string, error) {
	return "⏸️ **Music paused** (Mock response)", nil
}

func (u *MusicUsecase) handleResume(ctx context.Context) (string, error) {
	return "▶️ **Music resumed** (Mock response)", nil
}

func (u *MusicUsecase) handleQueue(ctx context.Context) (string, error) {
	return "📋 **Queue is empty** (Mock response)", nil
}

func (u *MusicUsecase) handleSkip(ctx context.Context) (string, error) {
	return "⏭️ **Song skipped** (Mock response)", nil
}

func (u *MusicUsecase) getHelpMessage() string {
	return `🎵 **Music Bot Commands**

▶️ **Playback:**
• ` + "`!play <song>`" + ` - Play a song
• ` + "`!pause`" + ` - Pause current song  
• ` + "`!resume`" + ` - Resume playback
• ` + "`!stop`" + ` - Stop and clear queue
• ` + "`!skip`" + ` - Skip current song

📋 **Queue:**
• ` + "`!queue`" + ` - Show current queue
• ` + "`!help`" + ` - Show this help

📝 **Examples:**
• ` + "`!play Never Gonna Give You Up`" + `
• ` + "`!play https://youtube.com/watch?v=...`" + `

⚠️ **Status**: Mock implementation - Full functionality coming soon!`
}
