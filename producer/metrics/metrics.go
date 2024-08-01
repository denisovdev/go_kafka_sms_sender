package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_messages",
		},
		[]string{"status"},
	)

	ProcessedMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "processed_messages",
		},
		[]string{},
	)
)

type Mertics struct {
	TotalMessages     *prometheus.CounterVec
	ProcessedMessages *prometheus.CounterVec
}

func MonitorTotalMessages(status string) {
	TotalMessages.WithLabelValues(status).Inc()
}

func MonitorProcessedMessages() {
	ProcessedMessages.WithLabelValues().Inc()
}
