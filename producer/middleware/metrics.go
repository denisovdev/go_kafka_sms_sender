package middleware

import (
	"fmt"

	"github.com/denisovdev/go_kafka_sms_sender/producer/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics(c *gin.Context) {
	c.Next()
	metrics.MonitorTotalMessages(fmt.Sprintf("%d", c.Writer.Status()))
}
