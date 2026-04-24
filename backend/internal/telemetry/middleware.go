// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UnmatchedRouteLabel is the bounded sentinel used as the route label for
// requests that did not match any registered chi pattern (typically 404s).
// Using the raw URL path here would blow Prometheus label cardinality for any
// attacker or confused client spraying arbitrary URLs, so we collapse all
// unmatched requests into one series.
const UnmatchedRouteLabel = "unmatched"

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
				request = request.WithContext(ctx)
			}

			if options.Metrics != nil {
				options.Metrics.InflightBegin()
			}

			// Finalise span/metrics in a defer so a panicking handler still
			// releases the in-flight gauge and records a request total.
			// chi.Recoverer should sit *inside* this middleware in the chain
			// (registered after it on the router) so it catches the panic,
			// calls WriteHeader(500) through statusWriter, and returns
			// normally — at which point this defer fires with the 500 status
			// already captured. If a panic ever escapes past this middleware
			// the defer still runs during unwind, leaving the in-flight gauge
			// and counter consistent (status defaults to the last value
			// statusWriter saw, typically 200).
			defer func() {
				route := resolveRoute(request)
				duration := time.Since(start).Seconds()

				if span != nil {
					span.SetAttributes(
						attribute.String("http.method", request.Method),
						attribute.String("http.route", route),
						attribute.Int("http.status_code", spanWriter.status),
					)
					span.SetName(request.Method + " " + route)
					span.End()
				}

				if options.Metrics != nil {
					options.Metrics.InflightEnd()
					options.Metrics.RecordHTTPRequest(request.Method, route, spanWriter.status, duration)
				}
			}()

			next.ServeHTTP(spanWriter, request)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (s *statusWriter) WriteHeader(code int) {
	if s.wroteHeader {
		return
	}
	s.status = code
	s.wroteHeader = true
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusWriter) Write(payload []byte) (int, error) {
	if !s.wroteHeader {
		s.wroteHeader = true
		// status already defaults to 200
	}
	return s.ResponseWriter.Write(payload)
}

// Unwrap lets net/http's ResponseController reach the underlying writer so
// handlers that need http.Hijacker, http.Flusher, or http.Pusher keep working
// behind the middleware.
func (s *statusWriter) Unwrap() http.ResponseWriter { return s.ResponseWriter }

func resolveRoute(request *http.Request) string {
	if ctx := chi.RouteContext(request.Context()); ctx != nil {
		if pattern := ctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return UnmatchedRouteLabel
}
