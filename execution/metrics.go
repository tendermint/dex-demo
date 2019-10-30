package execution

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ProcessingTime  metrics.Histogram
	OrdersProcessed metrics.Histogram
}

var exMetrics = &Metrics{
	ProcessingTime: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "dex",
		Subsystem: "execution",
		Name:      "execution_time",
		Help:      "Time for all match, and fill operations to complete.",
		Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
	}, []string{}),
	OrdersProcessed: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "dex",
		Subsystem: "execution",
		Name:      "orders_processed",
		Help:      "Number of orders processed.",
		Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
	}, []string{}),
}

func PrometheusMetrics() *Metrics {
	return exMetrics
}
