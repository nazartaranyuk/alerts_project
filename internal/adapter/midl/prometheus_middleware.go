package midl

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
}

func AddPrometheusMiddleware(server *echo.Echo) {
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			method := c.Request().Method

			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			status := c.Response().Status
			httpRequests.WithLabelValues(path, method, fmt.Sprintf("%d", status)).Inc()
			httpDuration.WithLabelValues(path, method, fmt.Sprintf("%d", status)).Observe(duration)

			return err
		}
	})
}
