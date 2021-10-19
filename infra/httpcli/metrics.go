package httpcli

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http/httpproxy"
)

func init() {
	prometheus.MustRegister(counter, histvec, inflight)
}

var (
	deftransp = &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	inflight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_client_in_flight_requests",
			Help: "Number of in-flight requests.",
		},
	)

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_requests_total",
			Help: "A counter for requests.",
		},
		[]string{"code", "method"},
	)

	histvec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_client_request_duration_seconds",
			Help:    "A histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)
)

// defaultTransport returns a RoundTripper with prometheus counters embed.
func defaultTransport() http.RoundTripper {
	return promhttp.InstrumentRoundTripperInFlight(
		inflight, promhttp.InstrumentRoundTripperCounter(
			counter, promhttp.InstrumentRoundTripperDuration(
				histvec, deftransp,
			),
		),
	)
}
