package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	Port string
	Mode string
}

type Processor struct {
	Topic            string
	TakeMessageLimit int
	ReservationTime  time.Duration
}

type Postgres struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSL_Mode string
}

type Producer struct {
	URL []string
}

type Config struct {
	App       *App
	Processor *Processor
	Postgres  *Postgres
	Producer  *Producer
}

var (
	ErrApp       = errors.New("invalid app config")
	ErrProcessor = errors.New("invalid processor config")
	ErrPostgres  = errors.New("invalid postgres config")
	ErrKafka     = errors.New("invalid kafka config")
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	app, err := newAppConfig()
	if err != nil {
		return nil, err
	}

	processor, err := newProcessorConfig()
	if err != nil {
		return nil, err
	}

	postgres, err := newPostgresConfig()
	if err != nil {
		return nil, err
	}

	producer, err := newProducerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		App:       app,
		Processor: processor,
		Postgres:  postgres,
		Producer:  producer,
	}, nil
}

func newAppConfig() (*App, error) {
	app := App{
		Port: os.Getenv("APP_PORT"),
		Mode: os.Getenv("APP_MODE"),
	}
	if app.Mode == "" || app.Port == "" {
		return nil, ErrApp
	}
	return &app, nil
}

func newProcessorConfig() (*Processor, error) {
	takeMessageLimit, err := strconv.Atoi(os.Getenv("PROCESSOR_TAKE_MESSAGE_LIMIT"))
	if err != nil {
		return nil, ErrProcessor
	}
	reservationTime, err := strconv.Atoi(os.Getenv("PROCESSOR_RESERVATION_TIME"))
	if err != nil {
		return nil, ErrProcessor
	}
	processor := Processor{
		Topic:            os.Getenv("PROCESSOR_TOPIC"),
		TakeMessageLimit: takeMessageLimit,
		ReservationTime:  time.Second * time.Duration(reservationTime),
	}
	if processor.Topic == "" {
		return nil, ErrProcessor
	}
	return &processor, nil
}

func newPostgresConfig() (*Postgres, error) {
	postgres := Postgres{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSL_Mode: os.Getenv("DB_SSL_MODE"),
	}
	if postgres.Host == "" || postgres.Port == "" || postgres.Database == "" || postgres.User == "" || postgres.Password == "" || postgres.SSL_Mode == "" {
		return nil, ErrPostgres
	}
	return &postgres, nil
}

func newProducerConfig() (*Producer, error) {
	producer := Producer{
		URL: strings.Fields(os.Getenv("PRODUCER_URL")),
	}
	if len(producer.URL) == 0 {
		return nil, ErrKafka
	}
	return &producer, nil
}
