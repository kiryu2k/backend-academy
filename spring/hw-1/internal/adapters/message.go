package adapters

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"ws-chat/internal/domain"
)

type messageRepo struct {
	pool *pgxpool.Pool
}

func NewMessageRepo(pool *pgxpool.Pool) messageRepo {
	return messageRepo{
		pool: pool,
	}
}

const (
	saveMessageQuery = `insert into messages (author, text, send_time) values ($1, $2, $3)`
	getMessagesQuery = `select author, text, send_time from messages order by send_time desc limit $1`
)

const recentMessageCount = 10

func (m messageRepo) SaveMessage(ctx context.Context, message domain.Message) error {
	_, err := m.pool.Exec(ctx, saveMessageQuery, message.Author, message.Text, message.Time)
	if err != nil {
		return errors.WithMessage(err, "insert message")
	}
	return nil
}

func (m messageRepo) GetRecentMessages(ctx context.Context) ([]domain.Message, error) {
	messages := make([]domain.Message, 0)
	rows, err := m.pool.Query(ctx, getMessagesQuery, recentMessageCount)
	if err != nil {
		return nil, errors.WithMessage(err, "select messages")
	}
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.Author, &msg.Text, &msg.Time); err != nil {
			return nil, errors.WithMessage(err, "scan rows")
		}
		messages = append(messages, msg)
	}
	slices.Reverse(messages)
	return messages, nil
}
