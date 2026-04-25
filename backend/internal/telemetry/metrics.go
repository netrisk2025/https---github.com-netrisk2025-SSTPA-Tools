// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry  *prometheus.Registry
	requests  *prometheus.CounterVec
	durations *prometheus.HistogramVec
	inflight  prometheus.Gauge
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sstpa",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total SSTPA API HTTP requests labelled by method, route, and status.",
		},
		[]string{"method", "route", "status"},
	)
	durations := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sstpa",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "SSTPA API HTTP request duration histogram.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)
	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "sstpa",
		Subsystem: "http",
		Name:      "inflight_requests",
		Help:      "In-flight SSTPA API HTTP requests.",
	})

	registry.MustRegister(requests, durations, inflight)
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Metrics{registry: registry, requests: requests, durations: durations, inflight: inflight}
}

func (m *Metrics) Registry() *prometheus.Registry { return m.registry }

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry})
}

// RecordHTTPRequest records one completed HTTP request. The route argument MUST
// be a bounded-cardinality router pattern (e.g. "/api/v1/nodes/{hid}"), never a
// concrete URL path — the chi route pattern is the canonical source and is
// produced by the telemetry.Middleware wrapper in middleware.go.
func (m *Metrics) RecordHTTPRequest(method string, route string, status int, durationSeconds float64) {
	labels := prometheus.Labels{
		"method": method,
		"route":  route,
		"status": strconv.Itoa(status),
	}
	m.requests.With(labels).Inc()
	m.durations.With(labels).Observe(durationSeconds)
}

func (m *Metrics) InflightBegin() { m.inflight.Inc() }
func (m *Metrics) InflightEnd()   { m.inflight.Dec() }
