package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"ws-chat/internal/domain"
)

type handler struct {
	service  domain.UseCase
	upgrader websocket.Upgrader
	logger   *zap.Logger
}

func newHandler(service domain.UseCase, logger *zap.Logger) handler {
	return handler{
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Пропускаем любой запрос
			},
		},
		logger: logger,
	}
}

const usernameKey = "X-User-Name-Key"

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	username := r.Header.Get(usernameKey)
	if username == "" {
		h.logger.Info(fmt.Sprintf("empty '%s' header", usernameKey))
		return
	}
	err = h.service.Handle(r.Context(), username, newClient(conn))
	if err != nil && !errors.Is(err, domain.ErrConnectionClosed) {
		h.logger.Error(err.Error())
	} else {
		h.logger.Info(fmt.Sprintf("user '%s' closed connection", username))
	}
}
