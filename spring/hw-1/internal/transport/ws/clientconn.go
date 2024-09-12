package ws

import (
	"github.com/gorilla/websocket"
	"ws-chat/internal/domain"
)

type client struct {
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) client {
	return client{conn: conn}
}

func (c client) WriteMessage(data []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c client) ReadMessage() ([]byte, error) {
	msgType, msg, err := c.conn.ReadMessage()
	if msgType == websocket.CloseMessage {
		return nil, domain.ErrConnectionClosed
	}
	return msg, err
}
