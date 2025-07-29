package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Discord AI Tech News Bot API", "status": "running"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "bot": "online"})
	})

	r.POST("/webhook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
	})
}
