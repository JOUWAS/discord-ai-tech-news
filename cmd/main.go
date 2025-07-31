package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"discord-ai-tech-news/config"
	"discord-ai-tech-news/internal/bot"

	// News Bot imports
	newsHandler "discord-ai-tech-news/internal/handlers/news"
	newsRepo "discord-ai-tech-news/internal/repositories/news"
	newsService "discord-ai-tech-news/internal/services/news"
	newsUsecase "discord-ai-tech-news/internal/usecases/news"

	// HTTP handler
	httpHandler "discord-ai-tech-news/internal/handler/http"
)

func main() {
	cfg := config.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Build dependencies dari luar ke dalam
	newsRepository := newsRepo.NewNewsApiRepository()
	newsServiceImpl := newsService.NewExternalNewsService(newsRepository)
	newsUsecaseImpl := newsUsecase.NewNewsUsecase(newsServiceImpl)
	newsHandlerImpl := newsHandler.NewNewsHandler(newsUsecaseImpl)

	// Production
	newsBot := bot.NewDiscordBot(cfg.DiscordToken, newsHandlerImpl)
	defer newsBot.Close()

	// Start Gin HTTP server
	router := gin.Default()
	httpHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Web server started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start: %s", err)
		}
	}()

	log.Println("Bot is now running. Press CTRL+C to exit.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
