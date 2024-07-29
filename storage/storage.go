package storage

import (
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/models"
)

type Messager interface {
	AddMessage(msg *models.Message) error
}

type Producer interface {
	TakeMessage(reserved_to time.Time, limit int) ([]*models.MessageStorage, error)
	UpdateStatus(id int) error
}
