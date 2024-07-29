package services

import (
	"github.com/denisovdev/go_kafka_sms_sender/models"
	"github.com/denisovdev/go_kafka_sms_sender/storage"
)

type MessagerService struct {
	Storage storage.Messager
}

func NewMessager(storage storage.Messager) *MessagerService {
	return &MessagerService{
		Storage: storage,
	}
}

func (service *MessagerService) CreateMessage(request *models.MessageRequest) error {
	// code generating and convert request structure to filled message structure
	code := "7219"
	message := request.ConvertToMessage(code)

	// send message structure to database layer
	err := service.Storage.AddMessage(message)
	if err != nil {
		return err
	}

	return nil
}
