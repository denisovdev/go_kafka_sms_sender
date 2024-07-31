package sms

import "github.com/denisovdev/go_kafka_sms_sender/consumer/models"

type Client interface {
	Send(*models.Message)
}
