package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewPrometheusMetrics() gin.HandlerFunc {
	requestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"host", "url", "method", "code"},
	)
	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_second",
			Help:    "The HTTP request latencies in second.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"url", "method", "code"},
	)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)
		status := strconv.Itoa(c.Writer.Status())
		url := c.Request.URL.Path
		requestsTotal.WithLabelValues(c.Request.Host, url, c.Request.Method, status).Inc()
		requestDuration.WithLabelValues(url, c.Request.Method, status).Observe(elapsed)
	}
}
