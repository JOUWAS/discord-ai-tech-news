package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not set in the environment variables")
	}

	log.Println("TOKEN is set to:", token)

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.MessageCreate) {
		if r.Author.ID == s.State.User.ID {
			return
		}

		log.Printf("=== Message Debug Info ===")
		log.Printf("Author: %s (ID: %s)", r.Author.Username, r.Author.ID)
		log.Printf("Content: '%s' (length: %d)", r.Content, len(r.Content))
		log.Printf("Channel: %s", r.ChannelID)
		log.Printf("Message ID: %s", r.ID)
		log.Printf("========================")

		if r.Content == "hello" {
			_, err := s.ChannelMessageSend(r.ChannelID, "sini kau hanip!!!")
			if err != nil {
				log.Printf("error sending message: %v", err)
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	defer sess.Close()

	log.Println("Bot is now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc

	// port := os.Getenv("APP_PORT")

	// r := gin.Default()
	// r.GET("/", func(c *gin.Context) {
	// 	c.String(200, "Hello, World!")
	// })

	// r.Run(":" + port)
}
