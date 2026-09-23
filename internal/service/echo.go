package service

import "context"

type EchoService struct{}

func NewEchoService() *EchoService {
	return &EchoService{}
}

func (s *EchoService) HandleMessage(ctx context.Context, msg Message) (*Message, error) {
	if msg.Text == "" {
		return nil, nil
	}
	return &Message{ChatID: msg.ChatID, Text: msg.Text}, nil
}
