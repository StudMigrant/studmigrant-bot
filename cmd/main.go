package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/StudMigrant/studmigrant-bot/internal/app"
	"github.com/StudMigrant/studmigrant-bot/internal/config"
	"github.com/StudMigrant/studmigrant-bot/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.New(os.Getenv("LOG_LEVEL"))

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("load config", "error", err)
		os.Exit(1)
	}

	a, err := app.New(ctx, *cfg, log)
	if err != nil {
		log.Error("init app", "error", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		log.Error("run app", "error", err)
		os.Exit(1)
	}
}
