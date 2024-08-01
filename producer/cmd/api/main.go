package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
	"github.com/denisovdev/go_kafka_sms_sender/producer/queue/kafka"
	"github.com/denisovdev/go_kafka_sms_sender/producer/server"
	"github.com/denisovdev/go_kafka_sms_sender/producer/services"
	"github.com/denisovdev/go_kafka_sms_sender/producer/storage/postgres"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := postgres.NewStorage(config.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	producer, err := kafka.NewProducer(config.Producer)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	processor := services.NewProcessor(storage, producer, config.Processor)
	srv := server.NewHTTPServer(config.App, storage)

	donech := make(chan os.Signal, 1)
	signal.Notify(donech, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		processor.StartProcessMessages(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Run(ctx); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-donech
	fmt.Println()
	cancel()

	wg.Wait()
}
