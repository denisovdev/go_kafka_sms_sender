package server

import (
	"fmt"

	"github.com/denisovdev/go_kafka_sms_sender/config"
	"github.com/denisovdev/go_kafka_sms_sender/services"
	"github.com/denisovdev/go_kafka_sms_sender/storage"
	"github.com/gin-gonic/gin"
)

type server struct {
	MessagerService services.Messager
	Config          *config.App
}

func NewHTTPServer(config *config.App, storage storage.Messager) *server {
	return &server{
		MessagerService: services.NewMessager(storage),
		Config:          config,
	}
}

func (srv *server) Run() error {
	gin.SetMode(srv.Config.Mode)

	router := gin.Default()
	router.POST("api/", srv.handleCreateMessage)

	addr := fmt.Sprintf("127.0.0.1:%s", srv.Config.Port)
	fmt.Printf("starting server at %s\n", addr)
	return router.Run(addr)
}
