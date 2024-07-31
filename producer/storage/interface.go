package storage

import (
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
)

type Storage interface {
	AddMessage(msg *models.Message) error
	TakeMessage(reserved_to time.Time, limit int) ([]*models.MessageStorage, error)
	UpdateStatus(id int) error
}
