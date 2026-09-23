package service

import "context"

type Message struct {
	ChatID int64
	Text   string
}

type Service interface {
	HandleMessage(ctx context.Context, msg Message) (*Message, error)
}
