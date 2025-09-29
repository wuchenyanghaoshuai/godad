package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Latency of HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	httpRequestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of in-flight HTTP requests",
	}, []string{"method", "path"})
)

// MetricsMiddleware records Prometheus metrics for each HTTP request.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		rawPath := c.Request.URL.Path
		httpRequestsInFlight.WithLabelValues(method, rawPath).Inc()
		start := time.Now()

		c.Next()

		httpRequestsInFlight.WithLabelValues(method, rawPath).Dec()

		route := c.FullPath()
		if route == "" {
			route = rawPath
		}
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(method, route, status).Inc()
		httpRequestsDuration.WithLabelValues(method, route).Observe(duration)
	}
}
