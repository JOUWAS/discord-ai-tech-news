package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"discord-ai-tech-news/internal/repository"
	"discord-ai-tech-news/internal/service"

	"github.com/bwmarrin/discordgo"
)

// Legacy function - untuk backward compatibility
func HandleDiscordMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Message from %s: %s", m.Author.Username, m.Content)

	if m.Content == "hello" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hello! 👋 Saya adalah bot AI Tech News!")
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

type MessageUsecase struct {
	newsService service.NewsService
}

func NewMessageUsecase(newsService service.NewsService) *MessageUsecase {
	return &MessageUsecase{
		newsService: newsService,
	}
}

func (u *MessageUsecase) ProcessMessage(ctx context.Context, content string) (string, error) {
	content = strings.TrimSpace(content)
	command := strings.ToLower(content)

	switch command {
	case "news", "berita", "tech", "teknologi":
		return u.handleNewsRequest(ctx)
	case "hello", "hi", "halo":
		return "Hello! 👋 Saya adalah **AI Tech News Bot**\n\n🤖 Saya bisa membantu Anda mendapatkan berita teknologi terbaru!\n\n💡 Ketik `help` untuk melihat command yang tersedia.", nil
	case "help", "bantuan":
		return u.getHelpMessage(), nil
	case "ping":
		return "🏓 Pong! Bot sedang online dan siap melayani!", nil
	case "status":
		return "✅ **Status Bot**: Online dan berjalan normal\n🔄 **Service**: News API Ready\n⚡ **Response Time**: < 1s", nil
	default:
		// Check if it's a search command
		if strings.HasPrefix(command, "search ") || strings.HasPrefix(command, "cari ") {
			keyword := strings.TrimPrefix(content, "search ")
			keyword = strings.TrimPrefix(keyword, "cari ")
			keyword = strings.TrimSpace(keyword)
			if keyword != "" {
				return u.handleSearchRequest(ctx, keyword)
			}
		}
		return u.getUnknownCommandMessage(), nil
	}
}

func (u *MessageUsecase) handleNewsRequest(ctx context.Context) (string, error) {
	newsResponse, err := u.newsService.FetchTechNews(ctx)
	if err != nil {
		log.Printf("Error fetching news: %v", err)
		return "❌ **Maaf, terjadi kesalahan saat mengambil berita**\n\n🔄 Silakan coba lagi dalam beberapa saat.\n💡 Atau ketik `help` untuk melihat command lainnya.", err
	}

	if len(newsResponse.News) == 0 {
		return "📰 **Tech News Update**\n\n🔍 Tidak ada berita teknologi terbaru saat ini.\n🔄 Coba lagi nanti untuk update terbaru!", nil
	}

	return u.newsService.FormatNewsForDiscord(newsResponse.News), nil
}

func (u *MessageUsecase) handleSearchRequest(ctx context.Context, keyword string) (string, error) {
	log.Printf("🔍 DEBUG: User searching for: %s", keyword)

	// Call search function from news service
	searchResults, err := u.newsService.SearchNews(ctx, keyword)
	if err != nil {
		log.Printf("❌ ERROR: Search failed for '%s': %v", keyword, err)
		return fmt.Sprintf("❌ **Pencarian Gagal**\n\n🔍 Tidak dapat mencari berita untuk: **%s**\n\n🔄 Silakan coba lagi atau gunakan keyword yang berbeda.", keyword), err
	}

	if len(searchResults) == 0 {
		return fmt.Sprintf("🔍 **Hasil Pencarian: \"%s\"**\n\n❌ Tidak ditemukan berita yang relevan.\n\n💡 **Tips:**\n• Coba keyword yang lebih umum\n• Gunakan bahasa Inggris (misal: AI, blockchain, startup)\n• Atau ketik `news` untuk berita terbaru", keyword), nil
	}

	// Format results for Discord
	return u.formatSearchResults(keyword, searchResults), nil
}

func (u *MessageUsecase) formatSearchResults(keyword string, results []repository.News) string {
	var response strings.Builder

	response.WriteString(fmt.Sprintf("🔍 **Hasil Pencarian: \"%s\"**\n\n", keyword))
	response.WriteString(fmt.Sprintf("📊 Ditemukan **%d artikel** yang relevan:\n\n", len(results)))

	// Limit to 5 results for Discord message length
	maxResults := 5
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	for i, article := range results {
		// Calculate time ago
		timeAgo := u.newsService.TimeAgo(article.PublishedAt)

		response.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, article.Title))
		if article.Description != "" {
			// Limit description length
			desc := article.Description
			if len(desc) > 150 {
				desc = desc[:150] + "..."
			}
			response.WriteString(fmt.Sprintf("📄 %s\n", desc))
		}
		response.WriteString(fmt.Sprintf("🔗 [Baca Selengkapnya](%s)\n", article.URL))
		response.WriteString(fmt.Sprintf("📅 %s • 📰 %s\n\n", timeAgo, article.Source))
	}

	if len(results) == maxResults {
		response.WriteString(fmt.Sprintf("💡 **Tips**: Gunakan keyword yang lebih spesifik untuk hasil yang lebih akurat.\n"))
	}

	return response.String()
}

func (u *MessageUsecase) getHelpMessage() string {
	return `📋 **AI Tech News Bot - Command List**

🔥 **Main Commands:**
• ` + "`news`" + ` atau ` + "`berita`" + ` - Dapatkan berita teknologi terbaru
• ` + "`hello`" + ` atau ` + "`hi`" + ` - Sapa bot
• ` + "`help`" + ` atau ` + "`bantuan`" + ` - Tampilkan menu ini
• ` + "`ping`" + ` - Cek status koneksi bot
• ` + "`status`" + ` - Lihat status bot

🔍 **Search Commands** *(Aktif)*:
• ` + "`search <keyword>`" + ` - Cari berita berdasarkan kata kunci
• ` + "`cari <keyword>`" + ` - Pencarian dalam bahasa Indonesia

📝 **Contoh Pencarian:**
• ` + "`search AI`" + ` - Cari berita tentang AI
• ` + "`cari blockchain`" + ` - Cari berita blockchain
• ` + "`search startup`" + ` - Cari berita startup

---
🤖 **About**: Saya adalah bot yang menyediakan berita teknologi terbaru dari berbagai sumber terpercaya.
📡 **Sources**: TechCrunch, Wired, The Verge, dan lainnya.
⚡ **Update**: Real-time news feed`
}

func (u *MessageUsecase) getUnknownCommandMessage() string {
	return `❓ **Command tidak dikenal**

🤔 Maaf, saya tidak mengerti command tersebut.

💡 **Coba command ini:**
• ` + "`news`" + ` - Berita teknologi terbaru
• ` + "`hello`" + ` - Sapa bot
• ` + "`help`" + ` - Lihat semua command

📝 **Tips**: Pastikan ejaan command benar dan tanpa typo!`
}
