package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/denisovdev/go_kafka_sms_sender/producer/config"
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
	var srv server
	handler := gin.Default()
	handler.POST("api/", srv.handleCreateMessage)
	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: handler,
	}
	srv.messagerService = services.NewMessager(storage)
	srv.config = config

	return &srv
}

func (srv *server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		srv.stop()
		fmt.Println("server stopped")
	}()

	fmt.Println("run server")
	return srv.httpServer.ListenAndServe()
}

func (srv *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		fmt.Println("can't stop server.. retrying")
		srv.stop()
	}
}
