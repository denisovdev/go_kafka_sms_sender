package queue

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/denisovdev/go_kafka_sms_sender/config"
)

type kafka struct {
	Conn  sarama.SyncProducer
	Topic string
}

func NewKafka(cfg *config.Kafka) (*kafka, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(cfg.URL, config)
	if err != nil {
		return nil, err
	}

	return &kafka{Conn: conn}, nil
}

func (k *kafka) Produce(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := k.Conn.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", k.Topic, partition, offset)
	return nil
}
