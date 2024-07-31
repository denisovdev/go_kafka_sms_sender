package kafka

import (
	"fmt"

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

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func (p *producer) Close() error {
	err := p.conn.Close()
	if err != nil {
		fmt.Println("can't close kafka producer connection")
		return err
	}
	fmt.Println("kafka producer connection closed")
	return nil
}

// func NewConsumer(config *config.Consumer) {
// 	consumer, err := kafka.NewProducer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"group.id":          "foo",
// 		"auto.offset.reset": "latest",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = consumer.Subscribe(config.Topic, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func (c *consumer) Consume() {
// 	for {
// 		ev := consumer.Poll(100)
// 		switch e := ev.(type) {
// 		case *kafka.Message:
// 			fmt.Printf("consumed new message drom the queue: %s\n", string(e.Value))
// 		case *kafka.Error:
// 			fmt.Printf("%v\n", e)
// 		}
// 	}
// }
