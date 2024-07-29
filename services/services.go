package services

import "github.com/denisovdev/go_kafka_sms_sender/models"

type Messager interface {
	CreateMessage(*models.MessageRequest) error
}
