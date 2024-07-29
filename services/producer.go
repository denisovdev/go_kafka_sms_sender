package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/queue"
	"github.com/denisovdev/go_kafka_sms_sender/storage"
)

type Producer struct {
	Storage storage.Producer
	Queue   queue.Queue
	Config  *config.Producer
}

func NewProducer(storage storage.Producer, queue queue.Queue, cfg *config.Producer) *Producer {
	return &Producer{
		Storage: storage,
		Queue:   queue,
		Config:  cfg,
	}
}

func (producer *Producer) StartProcessEvents(done *chan os.Signal, timout time.Duration) {
	ticker := time.NewTicker(timout)
	for {
		select {
		case <-*done:
			fmt.Println("\npusher stopped")
			return
		case <-ticker.C:
			messages, err := producer.Storage.TakeMessage(time.Now().Add(producer.Config.ReservationTime), producer.Config.TakeMessageLimit)
			if err != nil {
				fmt.Printf("can't take messages from database: %v\n", err)
				continue
			}
			if len(messages) == 0 {
				fmt.Println("no new messages to produce")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), producer.Config.ReservationTime)

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
					err = producer.Queue.Produce(producer.Config.Topic, byte_message)
					if err != nil {
						fmt.Printf("can't send message to queue: %v\n", err)
						continue
					}
					err = producer.Storage.UpdateStatus(message.ID)
					if err != nil {
						fmt.Printf("can't update status: %v\n", err)
					}
				}
			}
			cancel()
		}
	}
}
