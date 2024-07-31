package services

import (
	"context"
	"log"

	"github.com/denisovdev/go_kafka_sms_sender/consumer/external/sms"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/models"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/queue"
)

type sender struct {
	consumer queue.Consumer
	smsCient sms.Client
	channel  chan *models.QueueMessage
}

func NewSender(smsClient sms.Client, consumer queue.Consumer) *sender {
	return &sender{
		smsCient: smsClient,
		consumer: consumer,
		channel:  make(chan *models.QueueMessage),
	}
}

func (s *sender) Start(ctx context.Context) {
	log.Println("sender started")
	donech := make(chan struct{}, 1)

	go s.consumer.Consume(ctx, s.channel, donech)

	for {
		select {
		case queueMessage := <-s.channel:
			log.Printf("have new message | offset:%d\n", queueMessage.Offset)
			message, err := queueMessage.ConvertToMessage()
			if err != nil {
				log.Printf("can't convert QueueMessage to Message: %v\n", err)
				continue
			}
			s.smsCient.Send(message)
		case <-ctx.Done():
			if len(s.channel) == 0 && len(donech) == 1 {
				log.Println("sender stopped")
				return
			}
		}
	}
}
