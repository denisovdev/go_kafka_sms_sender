package queue

import (
	"context"

	"github.com/denisovdev/go_kafka_sms_sender/consumer/models"
)

type Consumer interface {
	Consume(ctx context.Context, messagech chan *models.QueueMessage, donech chan struct{})
}
