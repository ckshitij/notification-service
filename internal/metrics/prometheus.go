package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// APIRequestsDuration is a histogram of API request durations.
	APIRequestsDuration *prometheus.HistogramVec
	// SQLQueryDuration is a histogram of SQL query durations.
	SQLQueryDuration *prometheus.HistogramVec
)

func init() {
	APIRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_requests_duration_miliseconds",
			Help:    "Duration of API requests in miliseconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	SQLQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sql_query_duration_miliseconds",
			Help:    "Duration of SQL queries in miliseconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_name"},
	)
}

func PromHandler() http.Handler {
	registry := prometheus.NewRegistry()

	// Go runtime metrics (GC, memory, goroutines)
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(
		collectors.ProcessCollectorOpts{},
	))

	registry.MustRegister(APIRequestsDuration)
	registry.MustRegister(SQLQueryDuration)

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
