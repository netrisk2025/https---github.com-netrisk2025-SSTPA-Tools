// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracerProviderRecordsSpans(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)
	tracer := provider.Tracer("test")

	_, span := tracer.Start(context.Background(), "unit.span")
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Name() != "unit.span" {
		t.Fatalf("unexpected span name: %q", spans[0].Name())
	}
}

func TestShutdownIsIdempotent(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)

	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
}

func TestDisabledProviderShutdownIsNoop(t *testing.T) {
	provider, err := NewTracerProvider(context.Background(), TracerOptions{Enabled: false})
	if err != nil {
		t.Fatalf("disabled ctor: %v", err)
	}
	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}

	_, span := provider.Tracer("disabled").Start(context.Background(), "ignored")
	span.End()
	if span.SpanContext().IsSampled() {
		t.Fatal("disabled provider span should not be sampled")
	}
}
