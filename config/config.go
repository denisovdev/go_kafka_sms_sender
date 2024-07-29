package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Postgres struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSL_Mode string
}

type App struct {
	Port string
	Mode string
}

type Kafka struct {
	URL []string
}

type Producer struct {
	Topic            string
	TakeMessageLimit int
	ReservationTime  time.Duration
}
type config struct {
	Postgres *Postgres
	App      *App
	Kafka    *Kafka
	Producer *Producer
}

func NewConfig() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	app := App{
		Port: os.Getenv("APP_PORT"),
		Mode: os.Getenv("APP_MODE"),
	}
	if app.Mode == "" || app.Port == "" {
		return nil, errors.New("invalid app config")
	}

	takeMessageLimit, err := strconv.Atoi(os.Getenv("PRODUCER_TAKE_MESSAGE_LIMIT"))
	if err != nil {
		return nil, errors.New("invalid producer config")
	}
	reservationTime, err := strconv.Atoi(os.Getenv("PRODUCER_RESERVATION_TIME"))
	if err != nil {
		return nil, errors.New("invalid producer config")
	}
	producer := Producer{
		Topic:            os.Getenv("PRODUCER_TOPIC"),
		TakeMessageLimit: takeMessageLimit,
		ReservationTime:  time.Second * time.Duration(reservationTime),
	}
	if producer.Topic == "" {
		return nil, errors.New("invalid producer config")
	}

	postgres := Postgres{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSL_Mode: os.Getenv("DB_SSL_MODE"),
	}
	if postgres.Host == "" || postgres.Port == "" || postgres.Database == "" || postgres.User == "" || postgres.Password == "" || postgres.SSL_Mode == "" {
		return nil, errors.New("invalid postgres config")
	}

	kafka := Kafka{URL: strings.Fields(os.Getenv("KAFKA_URL"))}
	if len(kafka.URL) == 0 {
		return nil, errors.New("invalid kafka config")
	}

	return &config{
		Postgres: &postgres,
		App:      &app,
		Kafka:    &kafka,
		Producer: &producer,
	}, nil
}
