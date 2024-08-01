package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
)

type producer struct {
	conn sarama.SyncProducer
}

func NewProducer(cfg *config.Producer) (*producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(cfg.URL, config)
	if err != nil {
		return nil, err
	}

	log.Println("kafka producer connection open")
	return &producer{conn: conn}, nil
}

func (p *producer) Produce(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.conn.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func (p *producer) Close() error {
	err := p.conn.Close()
	if err != nil {
		log.Println("can't close kafka producer connection")
		return err
	}
	log.Println("kafka producer connection closed")
	return nil
}
