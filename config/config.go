package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	AppPort      string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	discordToken := os.Getenv("TOKEN")
	if discordToken == "" {
		log.Fatal("TOKEN is not set in environment variables")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &Config{
		DiscordToken: discordToken,
		AppPort:      appPort,
	}
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("TOKEN") == "" {
		log.Fatal("TOKEN is not set in environment variables")
	}
}
