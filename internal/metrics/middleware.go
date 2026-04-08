package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests proceeded by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors returned by the API",
	}, []string{"path", "status"})
)

func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status < 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		} else {
			HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		}
	}
}
