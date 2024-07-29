package postgres

import (
	"context"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type producerStorage struct {
	Pool *pgxpool.Pool
}

func NewProducerStorage(cfg *config.Postgres) (*producerStorage, error) {
	pool, err := newPool(cfg)
	if err != nil {
		return nil, err
	}

	return &producerStorage{Pool: pool}, nil
}

func (ps *producerStorage) TakeMessage(reserved_to time.Time, limit int) ([]*models.MessageStorage, error) {
	var messages []*models.MessageStorage
	query := `update message set "reserved_to" = $1 where "id" in (select "id" from message where "status" = 'new' and "reserved_to" < $2 limit $3) returning "id", "payload"`

	conn, err := ps.Pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), query, reserved_to, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message models.MessageStorage
		err := rows.Scan(&message.ID, &message.Payload)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (ps *producerStorage) UpdateStatus(id int) error {
	query := `update message set status = 'done' where id = $1`

	conn, err := ps.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
