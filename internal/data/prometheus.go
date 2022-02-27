package data

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MetricReqDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "finecloud",
		Subsystem:   "micro_service",
		Name:        "http_server_requests_seconds",
		Help:        "server requests duration(sec)",
		ConstLabels: nil,
		Buckets:     []float64{0.1, 0.25, 0.5, 0.75, 1},
	},
		[]string{"kind", "operation"})

	MetricReqTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "finecloud",
		Subsystem: "micro_service",
		Name:      "http_server_requests_seconds_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(MetricReqDurationHistogram, MetricReqTotalCounter)
}
