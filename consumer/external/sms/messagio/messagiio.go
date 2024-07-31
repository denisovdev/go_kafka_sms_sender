package messagio

import "github.com/denisovdev/go_kafka_sms_sender/consumer/models"

type messagio struct{}

func NewMessagioClient() *messagio {
	return &messagio{}
}

func (m *messagio) Send(message *models.Message) {
	// сборка и отправка http запроса с данными на внешний сервис отправки смс сообщений
}
