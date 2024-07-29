package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSL_Mode)
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, errors.New("unable to create connection pool")
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.New("unable to connect to database")
	}

	return pool, nil
}
