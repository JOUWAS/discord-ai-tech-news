package service

import (
	"context"
	"log"
	"net/http"
	"os"
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
	// Frankfurt is GMT+1 (GMT+2 during daylight saving)
	// To achieve GMT+7 timing, we need to subtract 6 hours (or 5 during daylight saving)
	// For safety, let's use 6 hours offset consistently

	// Auto news pada jam 08:00 WIB = 02:00 Frankfurt time
	_, err := cs.scheduler.NewJob(
		gocron.CronJob("0 2 * * *", false), // 02:00 Frankfurt = 08:00 WIB
		gocron.NewTask(cs.sendMorningNews),
	)
	if err != nil {
		return err
	}

	// Auto news pada jam 13:00 WIB = 07:00 Frankfurt time
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("0 7 * * *", false), // 07:00 Frankfurt = 13:00 WIB
		gocron.NewTask(cs.sendAfternoonNews),
	)
	if err != nil {
		return err
	}

	// Auto news pada jam 17:00 WIB = 11:00 Frankfurt time
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("0 11 * * *", false), // 11:00 Frankfurt = 17:00 WIB
		gocron.NewTask(cs.sendEveningNews),
	)
	if err != nil {
		return err
	}

	// Service health check - every minute
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("* * * * *", false), // Setiap menit
		gocron.NewTask(cs.pingStartEndpoint),
	)
	if err != nil {
		return err
	}

	// Test news job - 00:40 WIB = 18:40 Frankfurt time
	_, err = cs.scheduler.NewJob(
		gocron.CronJob("40 18 * * *", false), // 18:40 Frankfurt = 00:40 WIB
		gocron.NewTask(cs.sendTestNews),
	)
	if err != nil {
		return err
	}

	// Immediate test job - every minute for debugging
	// _, err = cs.scheduler.NewJob(
	// 	gocron.CronJob("* * * * *", false), // Every minute
	// 	gocron.NewTask(cs.sendImmediateTestNews),
	// )
	// if err != nil {
	// 	return err
	// }

	cs.scheduler.Start()
	log.Println("✅ Cron service started successfully")
	log.Println("📅 Auto news scheduled at: 08:00, 13:00, 17:00 WIB (02:00, 07:00, 11:00 Frankfurt)")
	log.Println("🧪 Test news scheduled at: 00:40 WIB (18:40 Frankfurt)")
	log.Println("🔄 Service health check: every minute")
	log.Println("🌍 Timezone: Adjusted for Frankfurt deployment to match GMT+7 (WIB)")
	return nil
}

func (cs *CronService) Stop() error {
	log.Println("🛑 Stopping cron service...")
	return cs.scheduler.Shutdown()
}

// Helper function to get current time in WIB (GMT+7)
func (cs *CronService) getWIBTime() time.Time {
	// Load WIB timezone (GMT+7)
	wibLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback: manually create GMT+7 offset
		wibLocation = time.FixedZone("WIB", 7*60*60) // 7 hours * 60 minutes * 60 seconds
	}
	return time.Now().In(wibLocation)
}

// Job untuk mengirim berita pagi (08:00 WIB)
func (cs *CronService) sendMorningNews() {
	wibTime := cs.getWIBTime()
	log.Printf("🌅 [AUTO NEWS] Sending morning tech news... (WIB: %s)", wibTime.Format("15:04"))
	cs.sendAutoNews("🌅 **Good Morning! Tech News Update**")
}

// Job untuk mengirim berita siang (13:00 WIB)
func (cs *CronService) sendAfternoonNews() {
	wibTime := cs.getWIBTime()
	log.Printf("🌞 [AUTO NEWS] Sending afternoon tech news... (WIB: %s)", wibTime.Format("15:04"))
	cs.sendAutoNews("🌞 **Afternoon Tech News Update**")
}

// Job untuk mengirim berita sore (17:00 WIB)
func (cs *CronService) sendEveningNews() {
	wibTime := cs.getWIBTime()
	log.Printf("🌆 [AUTO NEWS] Sending evening tech news... (WIB: %s)", wibTime.Format("15:04"))
	cs.sendAutoNews("🌆 **Evening Tech News Update**")
}

// Job untuk mengirim berita test (00:40 WIB)
func (cs *CronService) sendTestNews() {
	wibTime := cs.getWIBTime()
	log.Printf("🧪 [TEST NEWS] Sending test tech news... (WIB: %s)", wibTime.Format("15:04"))
	cs.sendAutoNews("🧪 **Test News Update - 00:40 WIB**")
}

// Job untuk mengirim berita test immediate (setiap menit untuk debugging)
// func (cs *CronService) sendImmediateTestNews() {
// 	wibTime := cs.getWIBTime()
// 	log.Printf("⚡ [IMMEDIATE TEST] Sending immediate test news... (WIB: %s)", wibTime.Format("15:04:05"))
// 	cs.sendAutoNews("⚡ **Immediate Test News - Every Minute**")
// }

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

	// Use WIB time for the timestamp
	wibTime := cs.getWIBTime()
	message := header + "\n\n" + formattedNews
	message += "\n\n---\n🤖 *Auto News Update* • " + wibTime.Format("15:04 WIB")

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

// Ping start endpoint setiap menit untuk menjaga service tetap aktif
func (cs *CronService) pingStartEndpoint() {
	// Get server URL from environment or use default
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080" // Default local server
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make POST request to /start endpoint
	resp, err := client.Post(serverURL+"/start", "application/json", nil)
	if err != nil {
		wibTime := cs.getWIBTime()
		log.Printf("⚠️ [HEALTH CHECK] Failed to ping /start endpoint: %v (WIB: %s)", err, wibTime.Format("15:04:05"))
		return
	}
	defer resp.Body.Close()

	// Check response status
	wibTime := cs.getWIBTime()
	if resp.StatusCode == http.StatusOK {
		log.Printf("✅ [HEALTH CHECK] Successfully pinged /start endpoint (WIB: %s)", wibTime.Format("15:04:05"))
	} else {
		log.Printf("⚠️ [HEALTH CHECK] /start endpoint returned status %d (WIB: %s)", resp.StatusCode, wibTime.Format("15:04:05"))
	}
}
