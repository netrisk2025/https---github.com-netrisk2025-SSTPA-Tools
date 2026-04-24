// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type TracerProvider struct {
	provider trace.TracerProvider
	sdk      *sdktrace.TracerProvider
	mu       sync.Mutex
	closed   bool
}

type TracerOptions struct {
	Enabled      bool
	OTLPEndpoint string
	ServiceName  string
}

func NewTracerProvider(ctx context.Context, options TracerOptions) (*TracerProvider, error) {
	if !options.Enabled {
		return &TracerProvider{provider: noop.NewTracerProvider()}, nil
	}

	endpoint, path, secure, err := parseOTLPEndpoint(options.OTLPEndpoint)
	if err != nil {
		return nil, err
	}

	exporterOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	if path != "" {
		exporterOpts = append(exporterOpts, otlptracehttp.WithURLPath(path))
	}
	if !secure {
		exporterOpts = append(exporterOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, exporterOpts...)
	if err != nil {
		return nil, err
	}

	sdkProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(buildResource(options.ServiceName)),
	)
	return &TracerProvider{provider: sdkProvider, sdk: sdkProvider}, nil
}

func NewTestTracerProvider(recorder *tracetest.SpanRecorder) *TracerProvider {
	sdkProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	return &TracerProvider{provider: sdkProvider, sdk: sdkProvider}
}

func (p *TracerProvider) Tracer(name string) trace.Tracer {
	return p.provider.Tracer(name)
}

func (p *TracerProvider) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	if p.sdk == nil {
		p.closed = true
		return nil
	}
	if err := p.sdk.Shutdown(ctx); err != nil {
		return err
	}
	p.closed = true
	return nil
}

func parseOTLPEndpoint(raw string) (string, string, bool, error) {
	if raw == "" {
		return "", "", false, errors.New("OTLP endpoint is empty")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", false, err
	}
	if parsed.Host == "" {
		return "", "", false, errors.New("OTLP endpoint missing host")
	}

	secure := strings.EqualFold(parsed.Scheme, "https")
	path := parsed.Path
	return parsed.Host, path, secure, nil
}

func buildResource(serviceName string) *sdkresource.Resource {
	if serviceName == "" {
		serviceName = "sstpa-backend"
	}
	return sdkresource.NewSchemaless(
		semconv.ServiceNameKey.String(serviceName),
	)
}
