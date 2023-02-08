package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var KafkaMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_messages",
		Help: "Number of messages from the topic",
	},
	[]string{"topic"},
)

var KafkaValidMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_valid_messages",
		Help: "Number of valid JSON messages from the topic",
	},
	[]string{"topic"},
)

var KafkaInvalidMessages = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kafka_invalid_messages",
		Help: "Number of invalid JSON messages from the topic",
	},
	[]string{"topic"},
)
