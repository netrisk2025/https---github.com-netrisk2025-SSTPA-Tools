// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestMiddlewareRecordsSpanAndMetric(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)
	metrics := NewMetrics()

	router := chi.NewRouter()
	router.Use(Middleware(MiddlewareOptions{Tracer: provider.Tracer("test"), Metrics: metrics}))
	router.Get("/api/v1/ping/{id}", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping/42", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", res.Code)
	}

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if !strings.Contains(spans[0].Name(), "/api/v1/ping/{id}") {
		t.Fatalf("span name should contain the route pattern, got %q", spans[0].Name())
	}

	metricsRecorder := httptest.NewRecorder()
	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler().ServeHTTP(metricsRecorder, metricsReq)
	body := metricsRecorder.Body.String()
	if !strings.Contains(body, `route="/api/v1/ping/{id}"`) {
		t.Fatalf("expected metrics to use chi route pattern, got:\n%s", body)
	}
}
