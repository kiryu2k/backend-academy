package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"ws-chat/internal/domain"
)

type hub struct {
	repo    domain.Repository
	logger  *zap.Logger
	clients map[string]domain.Client
	mu      *sync.Mutex
}

func New(repo domain.Repository, logger *zap.Logger) hub {
	return hub{
		repo:    repo,
		logger:  logger,
		clients: make(map[string]domain.Client),
		mu:      &sync.Mutex{},
	}
}

func (h hub) Handle(ctx context.Context, clientName string, client domain.Client) error {
	if err := h.addClient(clientName, client); err != nil {
		return errors.WithMessage(err, "add client")
	}
	defer h.removeClient(clientName)
	h.sendRecentMessages(ctx, client)
	for {
		msg, err := client.ReadMessage()
		if err != nil {
			break
		}
		if len(msg) == 0 {
			continue
		}
		h.logger.Info(string(msg), zap.String("client", clientName))
		err = h.repo.SaveMessage(ctx, domain.Message{
			Author: clientName,
			Text:   string(msg),
			Time:   time.Now(),
		})
		if err != nil {
			return errors.WithMessage(err, "save message")
		}
		go h.writeMessage(clientName, msg)
	}
	return nil
}

func (h hub) writeMessage(clientName string, message []byte) {
	msg := append([]byte(clientName+": "), message...)
	for _, client := range h.clients {
		err := client.WriteMessage(msg)
		if err != nil {
			h.logger.Warn(err.Error())
		}
	}
}

func (h hub) removeClient(clientName string) {
	h.mu.Lock()
	delete(h.clients, clientName)
	h.mu.Unlock()
}

func (h hub) addClient(clientName string, client domain.Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[clientName]; ok {
		h.logger.Info(fmt.Sprintf("client with name '%s' is already in chat", clientName))
		return errors.New("client with such name is already in chat")
	}
	h.clients[clientName] = client
	return nil
}

func (h hub) sendRecentMessages(ctx context.Context, client domain.Client) {
	recentMessages, err := h.repo.GetRecentMessages(ctx)
	if err != nil {
		h.logger.Warn(err.Error())
		return
	}
	for _, msg := range recentMessages {
		message := fmt.Sprintf("%s: %s", msg.Author, msg.Text)
		err := client.WriteMessage([]byte(message))
		if err != nil {
			h.logger.Warn(err.Error())
		}
	}
}
