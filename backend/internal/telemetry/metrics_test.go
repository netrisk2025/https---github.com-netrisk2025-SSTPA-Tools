// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsHandlerExposesSSTPAInstruments(t *testing.T) {
	metrics := NewMetrics()
	metrics.RecordHTTPRequest("GET", "/api/v1/nodes", 200, 0.012)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/metrics", nil)
	metrics.Handler().ServeHTTP(recorder, request)

	body := recorder.Body.String()
	for _, want := range []string{
		"sstpa_http_requests_total",
		"sstpa_http_request_duration_seconds_bucket",
		"go_goroutines",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected metrics output to include %q, got:\n%s", want, body)
		}
	}
}
