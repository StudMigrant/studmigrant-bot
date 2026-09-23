package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/StudMigrant/studmigrant-bot/internal/bot"
	"github.com/StudMigrant/studmigrant-bot/internal/config"
	"github.com/StudMigrant/studmigrant-bot/internal/handler"
	"github.com/StudMigrant/studmigrant-bot/internal/service"
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type App struct {
	cfg    config.Config
	bot    *bot.Bot
	server *http.Server
	log    *slog.Logger
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*App, error) {
	maxAPI, err := maxbot.NewApi(cfg.Token, maxbot.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}))
	if err != nil {
		return nil, fmt.Errorf("max api: %w", err)
	}

	if err := subscribeWebhook(ctx, maxAPI, cfg, log); err != nil {
		return nil, fmt.Errorf("subscribe webhook: %w", err)
	}

	svc := service.NewEchoService()
	b := bot.New(maxAPI, svc, log)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler.NewRouter(b.Webhook(cfg.WebhookSecret)),
	}

	return &App{cfg: cfg, bot: b, server: server, log: log}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		a.log.Info("server starting", "addr", a.cfg.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
		a.log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.log.Info("shutdown complete")
		return nil
	}
}

func subscribeWebhook(ctx context.Context, maxAPI *maxbot.Api, cfg config.Config, log *slog.Logger) error {
	subs, err := maxAPI.Subscriptions.GetSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("get subscriptions: %w", err)
	}

	for _, s := range subs.Subscriptions {
		if s.URL == cfg.WebhookURL {
			log.Info("webhook already subscribed", "url", cfg.WebhookURL)
			return nil
		}
	}

	updateTypes := []string{string(model.UpdateMessageCreated)}
	if _, err := maxAPI.Subscriptions.Subscribe(ctx, cfg.WebhookURL, cfg.WebhookSecret, updateTypes, ""); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	log.Info("webhook subscribed", "url", cfg.WebhookURL)
	return nil
}
