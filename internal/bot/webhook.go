package bot

import (
	"context"
	"net/http"

	"github.com/StudMigrant/studmigrant-bot/internal/service"
	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func (b *Bot) Webhook(secret string) http.Handler {
	return b.api.GetHandler(b.handle, secret)
}

func (b *Bot) handle(ctx context.Context, update model.Update) {
	switch update.UpdateType {
	case model.UpdateMessageCreated:
		b.onMessageCreated(ctx, update)
	}
}

func (b *Bot) onMessageCreated(ctx context.Context, update model.Update) {
	text := update.GetMessage().Body.Text
	if text == "" {
		return
	}

	b.log.Debug("message received", "chat_id", update.ChatID)

	out, err := b.svc.HandleMessage(ctx, service.Message{
		ChatID: update.ChatID,
		Text:   text,
	})
	if err != nil {
		b.log.Error("handle message", "chat_id", update.ChatID, "error", err)
		return
	}
	if out == nil {
		return
	}

	msg := maxbot.NewMessage().SetChat(out.ChatID).SetText(out.Text)
	if _, err := b.api.Messages.Send(ctx, msg); err != nil {
		b.log.Error("send message", "chat_id", out.ChatID, "error", err)
	}
}
