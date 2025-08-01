package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Discord AI Tech News Bot API", "status": "running"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "bot": "online"})
	})

	// Health check untuk cron jobs
	r.GET("/health/cron", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
			"cron_jobs": gin.H{
				"morning_news":   "08:00 daily",
				"afternoon_news": "13:00 daily",
				"evening_news":   "17:00 daily",
				"test_job":       "every 30 seconds",
			},
			"timezone":   "Asia/Jakarta (WIB)",
			"last_check": time.Now().Format("2006-01-02 15:04:05 WIB"),
		})
	})

	r.POST("/webhook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
	})
}
