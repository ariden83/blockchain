package metrics

import (
	"github.com/ariden83/blockchain/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	RouteCountReqs   *prometheus.CounterVec
	InFlight         prometheus.Gauge
	ResponseDuration *prometheus.HistogramVec
	RequestSize      *prometheus.SummaryVec
	ResponseSize     *prometheus.SummaryVec
	ApiParamsCounter *prometheus.CounterVec
}

func New(c config.Metrics) *Metrics {
	namespace := c.Namespace
	metric := &Metrics{
		RouteCountReqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_requests_total",
				Namespace:   namespace,
				Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
				ConstLabels: prometheus.Labels{"app": c.Name},
			},
			[]string{"code", "service"},
		),

		RequestSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        "http_push_size_bytes",
				Namespace:   namespace,
				Help:        "HTTP request size",
				Objectives:  map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
				ConstLabels: prometheus.Labels{"app": c.Name},
			},
			[]string{"service"},
		),

		ResponseSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        "http_response_size_bytes",
				Namespace:   namespace,
				Help:        "HTTP response size",
				Objectives:  map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
				ConstLabels: prometheus.Labels{"app": c.Name},
			},
			[]string{"service"},
		),

		ResponseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace:   namespace,
				Subsystem:   "http_request",
				Name:        "duration_seconds",
				Help:        "Run duration for each route",
				Buckets:     prometheus.ExponentialBuckets(0.01, 2, 25),
				ConstLabels: prometheus.Labels{"app": c.Name},
			},
			[]string{"service", "code"},
		),

		InFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "inflight_requests",
				Help:      "Number of HTTP requests currently processed",
			}),

		ApiParamsCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "count_params",
			Help:      "count the different paths used",
		}, []string{"limit", "nbOne", "nbTwo", "strOne", "strTwo"}),
	}

	prometheus.MustRegister(metric.RouteCountReqs)
	prometheus.MustRegister(metric.InFlight)
	prometheus.MustRegister(metric.ResponseDuration)
	prometheus.MustRegister(metric.ApiParamsCounter)
	return metric
}
