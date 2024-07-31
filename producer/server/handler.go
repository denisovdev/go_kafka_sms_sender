package server

import (
	"net/http"
	"regexp"

	"github.com/denisovdev/go_kafka_sms_sender/producer/models"
	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	Phone string `json:"phone" binding:"required,min=11"`
}

func (msgReq *MessageRequest) convertToMessage() *models.Message {
	return &models.Message{
		Phone: msgReq.Phone,
	}
}

func (msgReq *MessageRequest) validate() bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(msgReq.Phone)
}

func (srv *server) handleCreateMessage(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	// parse request body
	request := new(MessageRequest)
	err := ctx.Bind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	if !request.validate() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	message := request.convertToMessage()

	// send request structure to service layer
	err = srv.messagerService.CreateMessage(message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "message successfuly sended")
}
