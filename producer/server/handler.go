package server

import (
	"net/http"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
	"github.com/gin-gonic/gin"
)

func (srv *server) handleCreateMessage(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	// parse request body
	request := new(models.MessageRequest)
	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	if !request.Validate() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	message := request.ConvertToMessage()

	// send request structure to service layer
	err = srv.messagerService.CreateMessage(message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "message successfuly sent")
}
