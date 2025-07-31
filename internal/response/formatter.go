package response

import (
	"fmt"
	"strings"
)

// DiscordFormatter handles formatting responses for Discord display
type DiscordFormatter struct{}

// NewDiscordFormatter creates a new Discord formatter
func NewDiscordFormatter() *DiscordFormatter {
	return &DiscordFormatter{}
}

// FormatNewsResponse formats a NewsResponse for Discord display
func (f *DiscordFormatter) FormatNewsResponse(resp *NewsResponse) string {
	if !resp.Success {
		return f.formatError(resp.Error)
	}

	if len(resp.News) == 0 {
		return "📰 **Tech News Update**\n\n🔍 Tidak ada berita teknologi terbaru saat ini.\n🔄 Coba lagi nanti untuk update terbaru!"
	}

	var result strings.Builder
	result.WriteString("📰 **Tech News Update - Berita Teknologi Terbaru**\n\n")

	// Limit to 3 articles for Discord message length
	maxArticles := 3
	articles := resp.News
	if len(articles) > maxArticles {
		articles = articles[:maxArticles]
	}

	for i, article := range articles {
		result.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, article.Title))

		if article.Description != "" {
			description := article.Description
			if len(description) > 150 {
				description = description[:150] + "..."
			}
			result.WriteString(fmt.Sprintf("📝 %s\n", description))
		}

		result.WriteString(fmt.Sprintf("🔗 [Baca Selengkapnya](%s)\n", article.URL))
		result.WriteString(fmt.Sprintf("📅 %s • 📰 %s", article.TimeAgo, article.Source))

		// Add tags if available
		if len(article.Tags) > 0 {
			result.WriteString(fmt.Sprintf(" • 🏷️ %s", strings.Join(article.Tags[:min(3, len(article.Tags))], ", ")))
		}

		result.WriteString("\n\n")
	}

	if len(resp.News) > maxArticles {
		result.WriteString(fmt.Sprintf("📊 **Total**: %d artikel tersedia\n", len(resp.News)))
	}

	result.WriteString("---\n💡 *Ketik `help` untuk melihat command lainnya*")
	return result.String()
}

// FormatSearchResponse formats a SearchResponse for Discord display
func (f *DiscordFormatter) FormatSearchResponse(resp *SearchResponse) string {
	if !resp.Success {
		return f.formatError(resp.Error)
	}

	if len(resp.Results) == 0 {
		return fmt.Sprintf("🔍 **Hasil Pencarian: \"%s\"**\n\n❌ Tidak ditemukan berita yang relevan.\n\n💡 **Tips:**\n• Coba keyword yang lebih umum\n• Gunakan bahasa Inggris (misal: AI, blockchain, startup)\n• Atau ketik `news` untuk berita terbaru", resp.Query)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 **Hasil Pencarian: \"%s\"**\n\n", resp.Query))
	result.WriteString(fmt.Sprintf("📊 Ditemukan **%d artikel** yang relevan:\n\n", resp.ResultCount))

	// Limit to 5 results for Discord message length
	maxResults := 5
	results := resp.Results
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	for i, article := range results {
		result.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, article.Title))

		if article.Description != "" {
			description := article.Description
			if len(description) > 150 {
				description = description[:150] + "..."
			}
			result.WriteString(fmt.Sprintf("📄 %s\n", description))
		}

		result.WriteString(fmt.Sprintf("🔗 [Baca Selengkapnya](%s)\n", article.URL))
		result.WriteString(fmt.Sprintf("📅 %s • 📰 %s", article.TimeAgo, article.Source))

		// Add tags if available
		if len(article.Tags) > 0 {
			result.WriteString(fmt.Sprintf(" • 🏷️ %s", strings.Join(article.Tags[:min(3, len(article.Tags))], ", ")))
		}

		result.WriteString("\n\n")
	}

	if resp.ResultCount > maxResults {
		result.WriteString(fmt.Sprintf("💡 **Tips**: Gunakan keyword yang lebih spesifik untuk hasil yang lebih akurat. Total: %d artikel\n", resp.ResultCount))
	}

	return result.String()
}

// FormatBotResponse formats a BotResponse for Discord display
func (f *DiscordFormatter) FormatBotResponse(resp *BotResponse) string {
	if !resp.Success {
		return f.formatError(resp.Error)
	}

	if resp.DisplayText != "" {
		return resp.DisplayText
	}

	// Default formatting based on command
	switch resp.Command {
	case "hello", "hi", "halo":
		return "Hello! 👋 Saya adalah **AI Tech News Bot**\n\n🤖 Saya bisa membantu Anda mendapatkan berita teknologi terbaru!\n\n💡 Ketik `help` untuk melihat command yang tersedia."
	case "ping":
		return "🏓 Pong! Bot sedang online dan siap melayani!"
	case "help", "bantuan":
		return f.getHelpMessage()
	default:
		if resp.Message != "" {
			return resp.Message
		}
		return "✅ Perintah berhasil dijalankan."
	}
}

// FormatStatusResponse formats a StatusResponse for Discord display
func (f *DiscordFormatter) FormatStatusResponse(resp *StatusResponse) string {
	if !resp.Success {
		return f.formatError(resp.Error)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("✅ **Status Bot**: %s\n", resp.Status))

	if resp.Services != nil {
		result.WriteString("🔄 **Services**:\n")
		for service, status := range resp.Services {
			emoji := "✅"
			if status != "online" && status != "healthy" {
				emoji = "❌"
			}
			result.WriteString(fmt.Sprintf("  %s %s: %s\n", emoji, service, status))
		}
	}

	if resp.Performance != nil {
		result.WriteString("⚡ **Performance**:\n")
		if resp.Performance.ResponseTime != "" {
			result.WriteString(fmt.Sprintf("  📈 Response Time: %s\n", resp.Performance.ResponseTime))
		}
		if resp.Performance.MemoryUsage != "" {
			result.WriteString(fmt.Sprintf("  💾 Memory Usage: %s\n", resp.Performance.MemoryUsage))
		}
		if resp.Performance.ActiveUsers > 0 {
			result.WriteString(fmt.Sprintf("  👥 Active Users: %d\n", resp.Performance.ActiveUsers))
		}
	}

	return result.String()
}

// formatError formats error information for Discord display
func (f *DiscordFormatter) formatError(err *ErrorInfo) string {
	if err == nil {
		return "❌ **Terjadi kesalahan yang tidak diketahui**"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("❌ **Error**: %s\n", err.Message))

	if err.Details != "" {
		result.WriteString(fmt.Sprintf("📝 **Details**: %s\n", err.Details))
	}

	result.WriteString("\n🔄 Silakan coba lagi dalam beberapa saat.\n💡 Atau ketik `help` untuk melihat command lainnya.")

	return result.String()
}

// getHelpMessage returns the help message
func (f *DiscordFormatter) getHelpMessage() string {
	return `📋 **AI Tech News Bot - Command List**

🔥 **Main Commands:**
• ` + "`news`" + ` atau ` + "`berita`" + ` - Dapatkan berita teknologi terbaru
• ` + "`hello`" + ` atau ` + "`hi`" + ` - Sapa bot
• ` + "`help`" + ` atau ` + "`bantuan`" + ` - Tampilkan menu ini
• ` + "`ping`" + ` - Cek status koneksi bot
• ` + "`status`" + ` - Lihat status bot

🔍 **Search Commands**:
• ` + "`search <keyword>`" + ` - Cari berita berdasarkan kata kunci
• ` + "`cari <keyword>`" + ` - Pencarian dalam bahasa Indonesia

📝 **Contoh Pencarian:**
• ` + "`search AI`" + ` - Cari berita tentang AI
• ` + "`cari blockchain`" + ` - Cari berita blockchain
• ` + "`search startup`" + ` - Cari berita startup

---
🤖 **About**: Saya adalah bot yang menyediakan berita teknologi terbaru dari berbagai sumber terpercaya.
📡 **Sources**: Hacker News, TechCrunch, dan lainnya.
⚡ **Update**: Real-time news feed`
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
