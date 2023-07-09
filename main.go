package goprometheus

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	appHttpRequest = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_http_request_totals",
		Help: "The total number of application request http",
	}, []string{"method", "path"})

	appHttpCode = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_http_request_codes",
		Help: "The application request http status code",
	}, []string{"code"})

	appHttpLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_http_request_latency_seconds",
			Help: "Latency of HTTP requests.",
			// Define the desired histogram buckets.
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9},
		},
		[]string{"path"},
	)
)

func RecordHttpRequest(method string, path string) {
	appHttpRequest.WithLabelValues(method, path).Inc()
}

func RecordHttpCode(code int) {
	appHttpCode.WithLabelValues(strconv.Itoa(code)).Inc()
}

func RecordLatency(path string, start time.Time) {
	elapsed := time.Since(start).Seconds()
	appHttpLatency.WithLabelValues(path).Observe(elapsed)
}

type Middleware struct {
}

func (middleware *Middleware) MetricCollector() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			request := next(c)

			// Record metrics
			RecordHttpRequest(c.Request().Method, c.Path())
			RecordHttpCode(c.Response().Status)
			RecordLatency(c.Path(), start)

			return request
		}
	}
}
