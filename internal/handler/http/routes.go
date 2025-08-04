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
				"morning_news":   "08:00 WIB (02:00 Frankfurt)",
				"afternoon_news": "13:00 WIB (07:00 Frankfurt)",
				"evening_news":   "17:00 WIB (11:00 Frankfurt)",
				"test_news":      "00:25 WIB (18:25 Frankfurt)",
				"service_health": "every minute",
			},
			"deployment_timezone": "Frankfurt (GMT+1/GMT+2)",
			"target_timezone":     "Asia/Jakarta (GMT+7/WIB)",
			"timezone_note":       "Cron times adjusted for GMT+7 target from Frankfurt deployment",
			"last_check":          time.Now().Format("2006-01-02 15:04:05 MST"),
		})
	})

	r.POST("/webhook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
	})

	r.POST("/start", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "Service start triggered",
			"status":    "success",
			"timestamp": time.Now().Format("2006-01-02 15:04:05 WIB"),
		})
	})
}
