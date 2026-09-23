package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token         string
	WebhookURL    string
	WebhookSecret string
	Addr          string
}

func NewConfig() (*Config, error) {
	godotenv.Load()

	cfg := Config{
		Token:         os.Getenv("TOKEN"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		Addr:          os.Getenv("ADDR"),
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.Token == "" || cfg.WebhookURL == "" || cfg.WebhookSecret == "" {
		return nil, fmt.Errorf("MAX_BOT_TOKEN, WEBHOOK_URL, WEBHOOK_SECRET required")
	}
	return &cfg, nil
}
