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
	chimw "github.com/go-chi/chi/v5/middleware"
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

	body := scrapeMetrics(t, metrics)
	if !strings.Contains(body, `route="/api/v1/ping/{id}"`) {
		t.Fatalf("expected metrics to use chi route pattern, got:\n%s", body)
	}
}

func TestMiddlewareLabelsUnmatchedRequestsWithBoundedSentinel(t *testing.T) {
	metrics := NewMetrics()

	router := chi.NewRouter()
	router.Use(Middleware(MiddlewareOptions{Metrics: metrics}))
	// Register at least one route so chi wires up the middleware chain. Without
	// any routes, chi.Mux.ServeHTTP bypasses middlewares and calls the default
	// NotFoundHandler directly (mux.go handles the mx.handler == nil short path).
	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/no/such/path/abc123", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unmatched route, got %d", res.Code)
	}

	body := scrapeMetrics(t, metrics)
	if !strings.Contains(body, `route="unmatched"`) {
		t.Fatalf("expected unmatched sentinel label, got:\n%s", body)
	}
	if strings.Contains(body, `route="/no/such/path/abc123"`) {
		t.Fatalf("metrics leaked the raw URL path into a label: %s", body)
	}
}

func TestMiddlewareReleasesInflightGaugeOnPanic(t *testing.T) {
	metrics := NewMetrics()

	router := chi.NewRouter()
	// Order matters: the telemetry middleware must be OUTER so its defer runs
	// after chi.Recoverer has translated the panic into a 500 via WriteHeader.
	// chi wraps middlewares last-to-first, so middlewares registered later are
	// innermost. Registering Middleware first and Recoverer second yields:
	//   Request -> Middleware -> Recoverer -> Handler (panics)
	// Recoverer catches the panic, writes 500 through statusWriter, and
	// returns normally; then Middleware's defer records the 500 status and
	// releases the inflight gauge.
	router.Use(Middleware(MiddlewareOptions{Metrics: metrics}))
	router.Use(chimw.Recoverer)
	router.Get("/boom", func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 from chi.Recoverer, got %d", res.Code)
	}

	body := scrapeMetrics(t, metrics)
	if !strings.Contains(body, "sstpa_http_inflight_requests 0") {
		t.Fatalf("inflight gauge leaked after panic; metrics body:\n%s", body)
	}
	if !strings.Contains(body, `status="500"`) {
		t.Fatalf("expected a 500 sample in request totals; body:\n%s", body)
	}
}

func scrapeMetrics(t *testing.T, metrics *Metrics) string {
	t.Helper()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler().ServeHTTP(recorder, req)
	return recorder.Body.String()
}
