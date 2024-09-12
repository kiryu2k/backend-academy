package ws

import (
	"net/http"

	"go.uber.org/zap"
	"ws-chat/internal/domain"
)

func New(port string, service domain.UseCase, logger *zap.Logger) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: newHandler(service, logger),
	}
}
