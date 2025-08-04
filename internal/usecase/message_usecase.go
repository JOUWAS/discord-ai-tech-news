package usecase

import (
	"context"
 	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"discord-ai-tech-news/internal/response"
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
	formatter   *response.DiscordFormatter
}

func NewMessageUsecase(newsService service.NewsService) *MessageUsecase {
	return &MessageUsecase{
		newsService: newsService,
		formatter:   response.NewDiscordFormatter(),
	}
}

func (u *MessageUsecase) ProcessMessage(ctx context.Context, content string) (string, error) {
	content = strings.TrimSpace(content)
	originalCommand := strings.ToLower(content)

	botPrefixes := []string{"/", "!"}
	hasPrefix := false
	command := originalCommand

	for _, prefix := range botPrefixes {
		if strings.HasPrefix(originalCommand, prefix) {
			command = strings.TrimPrefix(originalCommand, prefix)
			hasPrefix = true
			break
		}
	}

	if !hasPrefix {
		return "", nil
	}

	switch command {
	case "news", "berita", "tech", "teknologi":
		return u.handleNewsRequest(ctx)
	case "hello", "hi", "halo", "hallo":
		resp := response.NewBotResponse("hello").
			WithDisplayText("Hello! 👋 Saya adalah **AI Tech News Bot Dev**\n\n🤖 Saya bisa membantu Anda mendapatkan berita teknologi terbaru!\n\n💡 Ketik `help` untuk melihat command yang tersedia.").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "help", "bantuan":
		resp := response.NewBotResponse("help").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "ping":
		resp := response.NewBotResponse("ping").
			WithDisplayText("🏓 Pong! Bot sedang online dan siap melayani!").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "status":
		services := map[string]string{
			"News API": "Ready",
			"Discord":  "Connected",
		}
		resp := response.NewStatusResponse().
			WithStatus("Online dan berjalan normal").
			WithServices(services).
			Build().(*response.StatusResponse)
		return u.formatter.FormatStatusResponse(resp), nil
	case "cron", "schedule", "jadwal":
		return u.handleCronStatusRequest(ctx)
	default:
		// Check if it's a search command
		if strings.HasPrefix(command, "search ") || strings.HasPrefix(command, "cari ") {
			keyword := strings.TrimPrefix(command, "search ")
			keyword = strings.TrimPrefix(keyword, "cari ")
			keyword = strings.TrimSpace(keyword)
			if keyword != "" {
				return u.handleSearchRequest(ctx, keyword)
			}
		}
		resp := response.NewBotResponse("unknown").
			WithDisplayText(u.getUnknownCommandMessage()).
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	}
}

func (u *MessageUsecase) handleNewsRequest(ctx context.Context) (string, error) {
	newsResponse, err := u.newsService.FetchTechNews(ctx)
	if err != nil {
		log.Printf("Error fetching news: %v", err)

		// Create error response
		errorResp := response.NewErrorResponse("NEWS_FETCH_ERROR", "Failed to fetch tech news").
			WithError("NEWS_FETCH_ERROR", "Maaf, terjadi kesalahan saat mengambil berita", err.Error()).
			Build().(*response.BaseResponse)

		return u.formatter.FormatBotResponse(&response.BotResponse{
			BaseResponse: *errorResp,
			Command:      "news",
		}), err
	}

	if len(newsResponse.News) == 0 {
		// Create empty response
		emptyResp := response.NewNewsResponse().
			WithMessage("No tech news available").
			Build().(*response.NewsResponse)

		return u.formatter.FormatNewsResponse(emptyResp), nil
	}

	// Create successful news response
	successResp := response.NewNewsResponse().
		WithNews(newsResponse.News).
		WithMessage("Latest tech news").
		Build().(*response.NewsResponse)

	return u.formatter.FormatNewsResponse(successResp), nil
}

func (u *MessageUsecase) handleSearchRequest(ctx context.Context, keyword string) (string, error) {
	log.Printf("🔍 DEBUG: User searching for: %s", keyword)

	// Call search function from news service
	searchResults, err := u.newsService.SearchNews(ctx, keyword)
	if err != nil {
		log.Printf("❌ ERROR: Search failed for '%s': %v", keyword, err)

		// Create error response
		errorResp := response.NewSearchResponse(keyword).
			WithError("SEARCH_ERROR", "Pencarian gagal", err.Error()).
			Build().(*response.SearchResponse)

		return u.formatter.FormatSearchResponse(errorResp), err
	}

	// Create search response
	searchResp := response.NewSearchResponse(keyword).
		WithSearchResults(searchResults, len(searchResults)).
		WithMessage("Search completed successfully").
		Build().(*response.SearchResponse)

	return u.formatter.FormatSearchResponse(searchResp), nil
}

func (u *MessageUsecase) handleCronStatusRequest(ctx context.Context) (string, error) {
	// Get the server URL from environment or use default
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080" // Default local server
	}

	// Make HTTP request to /health/cron endpoint
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(serverURL + "/health/cron")
	if err != nil {
		log.Printf("❌ ERROR: Failed to fetch cron status: %v", err)

		// Create error response
		errorResp := response.NewBotResponse("cron").
			WithDisplayText("❌ **Error**: Tidak dapat mengakses status cron jobs\n\n" +
				"🔧 **Possible Issues:**\n" +
				"• Server tidak berjalan\n" +
				"• Koneksi network bermasalah\n" +
				"• Endpoint `/health/cron` tidak tersedia").
			Build().(*response.BotResponse)

		return u.formatter.FormatBotResponse(errorResp), err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ ERROR: Failed to read response body: %v", err)

		errorResp := response.NewBotResponse("cron").
			WithDisplayText("❌ **Error**: Gagal membaca response dari server").
			Build().(*response.BotResponse)

		return u.formatter.FormatBotResponse(errorResp), err
	}

	// Parse JSON response
	var cronData map[string]interface{}
	if err := json.Unmarshal(body, &cronData); err != nil {
		log.Printf("❌ ERROR: Failed to parse JSON: %v", err)

		errorResp := response.NewBotResponse("cron").
			WithDisplayText("❌ **Error**: Gagal memparse response JSON dari server").
			Build().(*response.BotResponse)

		return u.formatter.FormatBotResponse(errorResp), err
	}

	// Build the response message
	var message strings.Builder
	message.WriteString("📅 **Cron Jobs Status**\n\n")

	// Status
	if status, ok := cronData["status"].(string); ok {
		message.WriteString(fmt.Sprintf("🔥 **Status**: %s\n\n", status))
	}

	// Cron Jobs
	if cronJobs, ok := cronData["cron_jobs"].(map[string]interface{}); ok {
		message.WriteString("⏰ **Scheduled Jobs:**\n")
		for jobName, schedule := range cronJobs {
			message.WriteString(fmt.Sprintf("• **%s**: %s\n", jobName, schedule))
		}
		message.WriteString("\n")
	}

	// Timezone
	if timezone, ok := cronData["timezone"].(string); ok {
		message.WriteString(fmt.Sprintf("🌍 **Timezone**: %s\n", timezone))
	}

	// Last Check
	if lastCheck, ok := cronData["last_check"].(string); ok {
		message.WriteString(fmt.Sprintf("🕐 **Last Check**: %s\n", lastCheck))
	}

	message.WriteString("\n💡 **Info**: Data diambil dari endpoint `/health/cron`")

	// Create successful response
	successResp := response.NewBotResponse("cron").
		WithDisplayText(message.String()).
		Build().(*response.BotResponse)

	return u.formatter.FormatBotResponse(successResp), nil
}

// ProcessMessageWithContext processes a message with user and channel context
func (u *MessageUsecase) ProcessMessageWithContext(ctx context.Context, content, userID, username, channelID, channelName string) (string, error) {
	content = strings.TrimSpace(content)
	originalCommand := strings.ToLower(content)

	botPrefixes := []string{"/", "!"}
	hasPrefix := false
	command := originalCommand

	for _, prefix := range botPrefixes {
		if strings.HasPrefix(originalCommand, prefix) {
			command = strings.TrimPrefix(originalCommand, prefix)
			hasPrefix = true
			break
		}
	}

	// If no prefix found, ignore the message
	if !hasPrefix {
		return "", nil // Return empty string to indicate message should be ignored
	}

	switch command {
	case "news", "berita", "tech", "teknologi":
		return u.handleNewsRequest(ctx)
	case "hello", "hi", "halo":
		resp := response.NewBotResponse("hello").
			WithDisplayText("Hello! 👋 Saya adalah **AI Tech News Bot**\n\n🤖 Saya bisa membantu Anda mendapatkan berita teknologi terbaru!\n\n💡 Ketik `help` untuk melihat command yang tersedia.").
			WithUserInfo(userID, username, false).
			WithChannelInfo(channelID, channelName, "text").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "help", "bantuan":
		resp := response.NewBotResponse("help").
			WithUserInfo(userID, username, false).
			WithChannelInfo(channelID, channelName, "text").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "ping":
		resp := response.NewBotResponse("ping").
			WithDisplayText("🏓 Pong! Bot sedang online dan siap melayani!").
			WithUserInfo(userID, username, false).
			WithChannelInfo(channelID, channelName, "text").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	case "status":
		services := map[string]string{
			"News API": "Ready",
			"Discord":  "Connected",
		}
		resp := response.NewStatusResponse().
			WithStatus("Online dan berjalan normal").
			WithServices(services).
			Build().(*response.StatusResponse)
		return u.formatter.FormatStatusResponse(resp), nil
	case "cron", "schedule", "jadwal":
		return u.handleCronStatusRequest(ctx)
	default:
		// Check if it's a search command
		if strings.HasPrefix(command, "search ") || strings.HasPrefix(command, "cari ") {
			keyword := strings.TrimPrefix(command, "search ")
			keyword = strings.TrimPrefix(keyword, "cari ")
			keyword = strings.TrimSpace(keyword)
			if keyword != "" {
				return u.handleSearchRequest(ctx, keyword)
			}
		}
		resp := response.NewBotResponse("unknown").
			WithDisplayText(u.getUnknownCommandMessage()).
			WithUserInfo(userID, username, false).
			WithChannelInfo(channelID, channelName, "text").
			Build().(*response.BotResponse)
		return u.formatter.FormatBotResponse(resp), nil
	}
}

func (u *MessageUsecase) getHelpMessage() string {
	return `📋 **AI Tech News Bot - Command List**

🔥 **Main Commands:**
• ` + "`/news`" + ` atau ` + "`!berita`" + ` - Dapatkan berita teknologi terbaru
• ` + "`/hello`" + ` atau ` + "`!hi`" + ` - Sapa bot
• ` + "`/help`" + ` atau ` + "`!bantuan`" + ` - Tampilkan menu ini
• ` + "`/ping`" + ` - Cek status koneksi bot
• ` + "`/status`" + ` - Lihat status bot
• ` + "`/cron`" + ` atau ` + "`!jadwal`" + ` - Lihat status cron jobs

🔍 **Search Commands**:
• ` + "`/search <keyword>`" + ` - Cari berita berdasarkan kata kunci
• ` + "`!cari <keyword>`" + ` - Pencarian dalam bahasa Indonesia

📝 **Contoh Penggunaan:**
• ` + "`/search AI`" + ` - Cari berita tentang AI
• ` + "`!cari blockchain`" + ` - Cari berita blockchain
• ` + "`/search startup`" + ` - Cari berita startup  
• ` + "`/cron`" + ` - Cek jadwal cron jobs

💡 **Tips**: Gunakan prefix ` + "`/`" + ` atau ` + "`!`" + ` di awal command

---
🤖 **About**: Saya adalah bot yang menyediakan berita teknologi terbaru dari berbagai sumber terpercaya.
📡 **Sources**: Hacker News dan sumber terpercaya lainnya.
⚡ **Update**: Real-time news feed`
}

func (u *MessageUsecase) getUnknownCommandMessage() string {
	return `❓ **Command tidak dikenal**

🤔 Maaf, saya tidak mengerti command tersebut.

💡 **Coba command ini:**
• ` + "`/news`" + ` - Berita teknologi terbaru
• ` + "`/hello`" + ` - Sapa bot
• ` + "`/help`" + ` - Lihat semua command

📝 **Tips**: 
• Gunakan prefix ` + "`/`" + ` atau ` + "`!`" + ` di awal command
• Pastikan ejaan command benar dan tanpa typo!`
}
