package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
	"github.com/denisovdev/go_kafka_sms_sender/producer/queue"
	"github.com/denisovdev/go_kafka_sms_sender/producer/storage"
)

type processor struct {
	Storage  storage.Storage
	Producer queue.Producer
	Config   *config.Processor
}

func NewProcessor(storage storage.Storage, producer queue.Producer, config *config.Processor) *processor {
	return &processor{
		Storage:  storage,
		Producer: producer,
		Config:   config,
	}
}

func (processor *processor) StartProcessMessages(ctx context.Context) {
	ticker := time.NewTicker(processor.Config.ReservationTime)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("processor stopped")
			return
		case <-ticker.C:
			messages, err := processor.Storage.TakeMessage(time.Now().Add(processor.Config.ReservationTime), processor.Config.TakeMessageLimit)
			if err != nil {
				fmt.Printf("can't take messages from database: %v\n", err)
				continue
			}
			if len(messages) == 0 {
				fmt.Println("no new messages to produce")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), processor.Config.ReservationTime)

		produce:
			for _, message := range messages {
				select {
				case <-ctx.Done():
					fmt.Println("done")
					break produce
				default:
					byte_message, err := json.Marshal(message.Payload)
					if err != nil {
						fmt.Printf("can't marshal message: %v\n", err)
						continue
					}
					err = processor.Producer.Produce(processor.Config.Topic, byte_message)
					if err != nil {
						fmt.Printf("can't send message to queue: %v\n", err)
						continue
					}
					err = processor.Storage.UpdateStatus(message.ID)
					if err != nil {
						fmt.Printf("can't update status: %v\n", err)
					}
				}
			}
			cancel()
		}
	}
}
