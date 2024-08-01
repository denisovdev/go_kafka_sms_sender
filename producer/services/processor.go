package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
	"github.com/denisovdev/go_kafka_sms_sender/producer/metrics"
	"github.com/denisovdev/go_kafka_sms_sender/producer/queue"
	"github.com/denisovdev/go_kafka_sms_sender/producer/storage"
)

type processor struct {
	storage  storage.Storage
	producer queue.Producer
	config   *config.Processor
}

func NewProcessor(storage storage.Storage, producer queue.Producer, config *config.Processor) *processor {
	return &processor{
		storage:  storage,
		producer: producer,
		config:   config,
	}
}

func (processor *processor) StartProcessMessages(ctx context.Context) {
	log.Println("processor started")
	ticker := time.NewTicker(processor.config.ReservationTime)
	for {
		select {
		case <-ctx.Done():
			log.Println("processor stopped")
			return
		case <-ticker.C:
			messages, err := processor.storage.TakeMessage(time.Now().Add(processor.config.ReservationTime), processor.config.TakeMessageLimit)
			if err != nil {
				log.Printf("can't take messages from database: %v\n", err)
				continue
			}
			if len(messages) == 0 {
				log.Println("no new messages to produce")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), processor.config.ReservationTime)

		produce:
			for _, message := range messages {
				select {
				case <-ctx.Done():
					log.Println("done")
					break produce
				default:
					byte_message, err := json.Marshal(message.Payload)
					if err != nil {
						log.Printf("can't marshal message: %v\n", err)
						continue
					}
					err = processor.producer.Produce(processor.config.Topic, byte_message)
					if err != nil {
						log.Printf("can't send message to queue: %v\n", err)
						continue
					}
					err = processor.storage.UpdateStatus(message.ID)
					if err != nil {
						log.Printf("can't update status: %v\n", err)
					}
					metrics.MonitorProcessedMessages()
				}
			}
			cancel()
		}
	}
}
