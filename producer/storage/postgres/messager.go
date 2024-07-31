package postgres

import (
	"context"
	"encoding/json"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
)

func (storage *Storage) AddMessage(message *models.Message) error {
	query := `insert into message (payload) values ($1)`
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn, err := storage.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), query, payload)
	return err
}
