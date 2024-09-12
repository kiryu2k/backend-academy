package domain

import (
	"context"
	"time"
)

type Message struct {
	Author string
	Text   string
	Time   time.Time
}

type Repository interface {
	SaveMessage(ctx context.Context, message Message) error
	GetRecentMessages(ctx context.Context) ([]Message, error)
}

type Client interface {
	WriteMessage(data []byte) error
	ReadMessage() ([]byte, error)
}

type UseCase interface {
	Handle(ctx context.Context, clientName string, client Client) error
}
