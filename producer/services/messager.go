package services

import (
	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
	"github.com/denisovdev/go_kafka_sms_sender/producer/storage"
	"github.com/denisovdev/go_kafka_sms_sender/producer/utils"
)

type MessagerService struct {
	Storage storage.Storage
}

func NewMessager(storage storage.Storage) *MessagerService {
	return &MessagerService{
		Storage: storage,
	}
}

func (service *MessagerService) CreateMessage(message *models.Message) error {
	// code generating and convert request structure to filled message structure
	code := utils.GenerateCode(6)
	message.Code = code

	// send message structure to database layer
	err := service.Storage.AddMessage(message)
	if err != nil {
		return err
	}

	return nil
}
