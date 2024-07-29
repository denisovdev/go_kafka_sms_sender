package main

import (
	"log"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/server"
	"github.com/denisovdev/go_kafka_sms_sender/storage/postgres"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	postgres, err := postgres.NewMessagerStorage(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewHTTPServer(cfg.App, postgres)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
