package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
	"github.com/denisovdev/go_kafka_sms_sender/producer/middleware"
	"github.com/denisovdev/go_kafka_sms_sender/producer/services"
	"github.com/denisovdev/go_kafka_sms_sender/producer/storage"
	"github.com/gin-gonic/gin"
)

type server struct {
	httpServer      *http.Server
	messagerService services.Messager
	config          *config.App
}

func NewHTTPServer(config *config.App, storage storage.Storage) *server {
	gin.SetMode(config.Mode)
	srv := new(server)
	handler := gin.Default()
	handler.POST("api/send/", middleware.Metrics, srv.handleCreateMessage)
	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: handler,
	}
	srv.messagerService = services.NewMessager(storage)
	srv.config = config

	return srv
}

func (srv *server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		srv.stop()
		log.Println("server stopped")
	}()

	log.Println("server started")
	return srv.httpServer.ListenAndServe()
}

func (srv *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		log.Println("can't stop server.. retrying")
		srv.stop()
	}
}
