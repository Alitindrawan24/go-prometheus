package goprometheus

import (
	"time"

	"github.com/labstack/echo/v4"
)

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
