package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/queue"
	"github.com/denisovdev/go_kafka_sms_sender/services"
	"github.com/denisovdev/go_kafka_sms_sender/storage/postgres"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	ProducerStorage, err := postgres.NewProducerStorage(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	kafka, err := queue.NewKafka(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}
	defer kafka.Conn.Close()

	producer := services.NewProducer(ProducerStorage, kafka, cfg.Producer)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	producer.StartProcessEvents(&done, time.Second*5)
}
