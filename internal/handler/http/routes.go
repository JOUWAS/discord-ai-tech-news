package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers HTTP routes
func RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Discord AI Tech News Bot is running",
			"bots": gin.H{
				"news":  "active",
				"music": "coming soon",
			},
		})
	})

	// Status endpoint
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Discord Multi-Bot",
			"version": "1.0.0",
			"bots": []gin.H{
				{
					"name":     "News Bot",
					"status":   "active",
					"features": []string{"news", "search", "help"},
				},
				{
					"name":     "Music Bot",
					"status":   "development",
					"features": []string{"play", "pause", "queue"},
				},
			},
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Discord AI Tech News Bot API",
			"endpoints": gin.H{
				"/health": "Health check",
				"/status": "Bot status information",
			},
		})
	})
}
