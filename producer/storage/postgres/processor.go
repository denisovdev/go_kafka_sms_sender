package postgres

import (
	"context"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/storage"
)

func (s *Storage) TakeMessage(reserved_to time.Time, limit int) ([]*storage.MessageStorage, error) {
	messages := make([]*storage.MessageStorage, 0, limit)
	query := `update message set "reserved_to" = $1 where "id" in (select "id" from message where "status" = 'new' and "reserved_to" < $2 order by created_at desc limit $3) returning "id", "payload"`

	conn, err := s.pool.Acquire(context.Background())
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
		var (
			id      int
			payload map[string]string
		)

		err := rows.Scan(&id, &payload)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &storage.MessageStorage{
			ID:      id,
			Payload: payload,
		})
	}
	return messages, nil
}

func (s *Storage) UpdateStatus(id int) error {
	query := `update message set status = 'done' where id = $1`

	conn, err := s.pool.Acquire(context.Background())
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
