package service

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type CronService struct {
	scheduler   gocron.Scheduler
	newsService NewsService
	discordBot  DiscordBotInterface // Tambahkan interface untuk Discord bot
}

// Interface untuk Discord bot
type DiscordBotInterface interface {
	SendNewsToChannel(channelName string, message string) error
}

func NewCronService(newsService NewsService, discordBot DiscordBotInterface) *CronService {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	return &CronService{
		scheduler:   scheduler,
		newsService: newsService,
		discordBot:  discordBot,
	}
}

func (cs *CronService) Start() error {
	// Auto news pada jam 08:00
	_, err := cs.scheduler.NewJob(
		gocron.CronJob("0 8 * * *", false), // Setiap hari jam 08:00
		gocron.NewTask(cs.sendMorningNews),
	)
	if err != nil {
		return err
	}

	// Auto news pada jam 13:00
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("0 13 * * *", false), // Setiap hari jam 13:00
		gocron.NewTask(cs.sendAfternoonNews),
	)
	if err != nil {
		return err
	}

	// Auto news pada jam 17:00
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("0 17 * * *", false), // Setiap hari jam 17:00
		gocron.NewTask(cs.sendEveningNews),
	)
	if err != nil {
		return err
	}

	cs.scheduler.Start()
	log.Println("✅ Cron service started successfully")
	log.Println("📅 Auto news scheduled at: 08:00, 13:00, 17:00")
	return nil
}

func (cs *CronService) Stop() error {
	log.Println("🛑 Stopping cron service...")
	return cs.scheduler.Shutdown()
}

// Job untuk mengirim berita pagi (08:00)
func (cs *CronService) sendMorningNews() {
	log.Println("🌅 [AUTO NEWS] Sending morning tech news...")
	cs.sendAutoNews("🌅 **Good Morning! Tech News Update**")
}

// Job untuk mengirim berita siang (13:00)
func (cs *CronService) sendAfternoonNews() {
	log.Println("🌞 [AUTO NEWS] Sending afternoon tech news...")
	cs.sendAutoNews("🌞 **Afternoon Tech News Update**")
}

// Job untuk mengirim berita sore (17:00)
func (cs *CronService) sendEveningNews() {
	log.Println("🌆 [AUTO NEWS] Sending evening tech news...")
	cs.sendAutoNews("🌆 **Evening Tech News Update**")
}

// Fungsi utama untuk mengambil dan mengirim berita
func (cs *CronService) sendAutoNews(header string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Ambil berita teknologi terbaru
	newsResponse, err := cs.newsService.FetchTechNews(ctx)
	if err != nil {
		log.Printf("❌ [AUTO NEWS] Error getting news: %v", err)
		// Kirim pesan error ke Discord
		errorMsg := "❌ **Tech News Update**\n\nMaaf, terjadi kesalahan saat mengambil berita teknologi terbaru. Silakan coba lagi nanti."
		cs.sendToDiscord(errorMsg)
		return
	}

	// Format pesan untuk Discord
	message := cs.formatNewsMessage(header, newsResponse)

	// Kirim ke Discord
	cs.sendToDiscord(message)

	log.Println("✅ [AUTO NEWS] News sent successfully to Discord")
}

// Format pesan berita untuk Discord
func (cs *CronService) formatNewsMessage(header string, newsResponse *NewsResponse) string {
	if newsResponse == nil || len(newsResponse.News) == 0 {
		return header + "\n\n❌ Tidak ada berita teknologi terbaru yang tersedia saat ini."
	}

	// Gunakan formatter yang sudah ada di NewsService
	formattedNews := cs.newsService.FormatNewsForDiscord(newsResponse.News)

	message := header + "\n\n" + formattedNews
	message += "\n\n---\n🤖 *Auto News Update* • " + time.Now().Format("15:04 WIB")

	return message
}

// Kirim pesan ke Discord channel
func (cs *CronService) sendToDiscord(message string) {
	// Channel names - coba beberapa kemungkinan format
	channelNames := []string{
		"🔥┃ai-tech-news", // Format dengan emoji separator
		"ai-tech-news",   // Format simple
		"tech-news",      // Format alternatif
		"general",        // Fallback ke general channel
	}

	var lastErr error
	for _, channelName := range channelNames {
		err := cs.discordBot.SendNewsToChannel(channelName, message)
		if err == nil {
			log.Printf("✅ [AUTO NEWS] Message sent to Discord channel '%s' successfully", channelName)
			return
		}
		lastErr = err
		log.Printf("⚠️ [AUTO NEWS] Failed to send to channel '%s': %v", channelName, err)
	}

	log.Printf("❌ [AUTO NEWS] Failed to send to any Discord channel: %v", lastErr)
}

// Hello World job untuk testing
func (cs *CronService) helloWorldJob() {
	log.Printf("👋 [HELLO WORLD] Hello World! - %s", time.Now().Format("15:04:05"))
}
