package storage

import (
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
)

type MessageStorage struct {
	ID      int               `json:"id"`
	Payload map[string]string `json:"payload"`
}
type Storage interface {
	AddMessage(msg *models.Message) error
	TakeMessage(reserved_to time.Time, limit int) ([]*MessageStorage, error)
	UpdateStatus(id int) error
}
