package postgres

import (
	"context"
	"encoding/json"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messagerStorage struct {
	Pool *pgxpool.Pool
}

func NewMessagerStorage(cfg *config.Postgres) (*messagerStorage, error) {
	pool, err := newPool(cfg)
	if err != nil {
		return nil, err
	}

	return &messagerStorage{Pool: pool}, nil
}

func (ms *messagerStorage) AddMessage(message *models.Message) error {
	query := `insert into message (payload) values ($1)`
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn, err := ms.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), query, payload)
	return err
}
