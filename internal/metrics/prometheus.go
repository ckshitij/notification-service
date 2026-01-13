package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PromHandler() http.Handler {
	registry := prometheus.NewRegistry()

	// Go runtime metrics (GC, memory, goroutines)
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(
		collectors.ProcessCollectorOpts{},
	))

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
