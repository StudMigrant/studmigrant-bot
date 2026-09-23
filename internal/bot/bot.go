package bot

import (
	"log/slog"

	"github.com/StudMigrant/studmigrant-bot/internal/service"
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

type Bot struct {
	api *maxbot.Api
	svc service.Service
	log *slog.Logger
}

func New(api *maxbot.Api, svc service.Service, log *slog.Logger) *Bot {
	return &Bot{api: api, svc: svc, log: log}
}
