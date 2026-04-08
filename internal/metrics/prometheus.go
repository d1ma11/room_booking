package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var customRegistry = prometheus.NewRegistry()

func init() {
	customRegistry.MustRegister(HttpRequestTotal, HttpRequestErrorTotal)
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
