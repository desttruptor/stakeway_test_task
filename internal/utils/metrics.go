package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "validator_api_requests_total",
			Help: "The total number of requests by endpoint",
		},
		[]string{"endpoint", "method", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "validator_api_request_duration_seconds",
			Help:    "The duration of requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	TasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "validator_api_tasks_total",
			Help: "The total number of async tasks by status",
		},
		[]string{"status"},
	)

	TaskDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "validator_api_task_duration_seconds",
			Help:    "The duration of async tasks in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)
