package ctrls

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_server_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"code", "handler", "method"},
	)
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Count of HTTP requests",
		},
		[]string{"code", "handler", "method"},
	)
)

func init() {
	prometheus.MustRegister(requestsDuration, requestsTotal)
}

// instrumentHandler wraps provided Handler with prometheus metrics for duration and total.
func instrumentHandler(name string, fn http.HandlerFunc) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		requestsDuration.MustCurryWith(prometheus.Labels{"handler": name}),
		promhttp.InstrumentHandlerCounter(
			requestsTotal.MustCurryWith(prometheus.Labels{"handler": name}),
			fn,
		),
	)
}
