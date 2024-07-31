package services

import (
	"os"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
)

type Messager interface {
	CreateMessage(*models.Message) error
}

type Processor interface {
	StartProcessMessages(done *chan os.Signal, timout time.Duration)
}
