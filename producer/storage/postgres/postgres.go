package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(config *config.Postgres) (*Storage, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Port, config.Database, config.SSL_Mode)
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, errors.New("unable to create connection pool")
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.New("unable to connect to database")
	}

	log.Println("postgres connection pool open")
	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
	log.Println("postgres connection closed")
}
