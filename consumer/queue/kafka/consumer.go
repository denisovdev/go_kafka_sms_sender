package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/config"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/models"
)

type consumer struct {
	conn *kafka.Consumer
}

func NewConsumer(config *config.Consumer) (*consumer, error) {
	conn, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.URL,
		"group.id":          config.GroupID,
		"auto.offset.reset": config.AutoOffsetReset,
	})
	if err != nil {
		return nil, err
	}

	err = conn.Subscribe(config.Topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("kafka consumer connection open")
	return &consumer{
		conn: conn,
	}, nil
}

func (c *consumer) Consume(ctx context.Context, messagech chan *models.QueueMessage, donech chan struct{}) {
	log.Println("kafka consumer started")
	for {
		select {
		case <-ctx.Done():
			log.Println("kafka consumer stopped")
			donech <- struct{}{}
			return
		default:
			event := c.conn.Poll(100)
			switch e := event.(type) {
			case *kafka.Message:
				queueMessage := &models.QueueMessage{
					Offset:  int(e.TopicPartition.Offset),
					Payload: e.Value,
				}
				messagech <- queueMessage
			case *kafka.Error:
				log.Printf("%v\n", e)
			}
		}
	}
}

func (c *consumer) Close() error {
	log.Println("kafka consumer connection closed")
	return c.conn.Close()
}
