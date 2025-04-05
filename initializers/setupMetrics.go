package initializers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_app_processed_ops_total",
		Help: "The total number of processed events",
	})
)
