// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package config

import (
	"os"
	"testing"
)

func TestLoadDefaultTelemetrySettings(t *testing.T) {
	os.Unsetenv("SSTPA_OTLP_ENDPOINT")
	os.Unsetenv("SSTPA_SERVICE_NAME")
	os.Unsetenv("SSTPA_TELEMETRY_METRICS")
	os.Unsetenv("SSTPA_TELEMETRY_TRACING")

	cfg := Load()
	if cfg.OTLPEndpoint != "http://otel-collector:4318" {
		t.Fatalf("OTLPEndpoint default = %q", cfg.OTLPEndpoint)
	}
	if cfg.ServiceName != "sstpa-backend" {
		t.Fatalf("ServiceName default = %q", cfg.ServiceName)
	}
	if !cfg.MetricsEnabled {
		t.Fatalf("MetricsEnabled default must be true")
	}
	if !cfg.TracingEnabled {
		t.Fatalf("TracingEnabled default must be true")
	}
}

func TestLoadOverridesTelemetrySettings(t *testing.T) {
	t.Setenv("SSTPA_OTLP_ENDPOINT", "http://localhost:4318")
	t.Setenv("SSTPA_SERVICE_NAME", "sstpa-test")
	t.Setenv("SSTPA_TELEMETRY_METRICS", "false")
	t.Setenv("SSTPA_TELEMETRY_TRACING", "false")

	cfg := Load()
	if cfg.OTLPEndpoint != "http://localhost:4318" {
		t.Fatalf("OTLPEndpoint override = %q", cfg.OTLPEndpoint)
	}
	if cfg.ServiceName != "sstpa-test" {
		t.Fatalf("ServiceName override = %q", cfg.ServiceName)
	}
	if cfg.MetricsEnabled {
		t.Fatalf("MetricsEnabled override must be false")
	}
	if cfg.TracingEnabled {
		t.Fatalf("TracingEnabled override must be false")
	}
}
