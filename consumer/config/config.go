package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrConsumer = errors.New("invalid consumer config")
)

type Consumer struct {
	URL             string
	Topic           string
	GroupID         string
	AutoOffsetReset string
}

type Config struct {
	Consumer *Consumer
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	consumer, err := newConsumerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Consumer: consumer,
	}, nil
}

func newConsumerConfig() (*Consumer, error) {
	consumer := Consumer{
		URL:             os.Getenv("CONSUMER_URL"),
		Topic:           os.Getenv("CONSUMER_TOPIC"),
		GroupID:         os.Getenv("CONSUMER_GROUP_ID"),
		AutoOffsetReset: os.Getenv("CONSUMER_AUTO_OFFSET_RESET"),
	}
	if consumer.URL == "" || consumer.Topic == "" || consumer.GroupID == "" || consumer.AutoOffsetReset == "" {
		return nil, ErrConsumer
	}
	return &consumer, nil
}
