// Package middleware Prometheus 监控指标中间件
package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "his_http_requests_total",
			Help: "HTTP 请求总数，按 method、path、status 分类",
		},
		[]string{"service", "method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "his_http_request_duration_seconds",
			Help:    "HTTP 请求处理耗时分布",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path"},
	)

	httpRequestInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "his_http_requests_in_flight",
			Help: "当前正在处理的 HTTP 请求数",
		},
		[]string{"service"},
	)
)

// Metrics Prometheus 监控指标中间件
func Metrics(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		httpRequestInFlight.WithLabelValues(serviceName).Inc()
		defer httpRequestInFlight.WithLabelValues(serviceName).Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(serviceName, c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(serviceName, c.Request.Method, path).Observe(duration)
	}
}
