package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/denisovdev/go_kafka_sms_sender/consumer/config"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/external/sms/messagio"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/queue/kafka"
	"github.com/denisovdev/go_kafka_sms_sender/consumer/services"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	messagio := messagio.NewMessagioClient()
	consumer, err := kafka.NewConsumer(config.Consumer)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	donech := make(chan os.Signal, 1)
	signal.Notify(donech, os.Interrupt, syscall.SIGTERM)

	sender := services.NewSender(messagio, consumer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sender.Start(ctx)
	}()

	<-donech
	fmt.Println()
	cancel()

	wg.Wait()
}
