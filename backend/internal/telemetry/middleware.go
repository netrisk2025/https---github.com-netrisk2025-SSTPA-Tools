// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareOptions struct {
	Tracer  trace.Tracer
	Metrics *Metrics
}

func Middleware(options MiddlewareOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			spanWriter := &statusWriter{ResponseWriter: writer, status: http.StatusOK}

			var (
				ctx  = request.Context()
				span trace.Span
			)
			if options.Tracer != nil {
				ctx, span = options.Tracer.Start(request.Context(), request.Method+" "+request.URL.Path)
				defer span.End()
				request = request.WithContext(ctx)
			}

			if options.Metrics != nil {
				options.Metrics.InflightBegin()
			}

			next.ServeHTTP(spanWriter, request)

			route := chiRoutePattern(request)
			duration := time.Since(start).Seconds()

			if span != nil {
				span.SetAttributes(
					attribute.String("http.method", request.Method),
					attribute.String("http.route", route),
					attribute.Int("http.status_code", spanWriter.status),
					attribute.String("http.status_text", strconv.Itoa(spanWriter.status)),
				)
				if route != "" {
					span.SetName(request.Method + " " + route)
				}
			}

			if options.Metrics != nil {
				options.Metrics.InflightEnd()
				target := route
				if target == "" {
					target = request.URL.Path
				}
				options.Metrics.RecordHTTPRequest(request.Method, target, spanWriter.status, duration)
			}
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func chiRoutePattern(request *http.Request) string {
	if ctx := chi.RouteContext(request.Context()); ctx != nil {
		if pattern := ctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return request.URL.Path
}
