package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const msMultiplier = 1000

var (
	requestsDuration *prometheus.HistogramVec
	requestsCounter  *prometheus.CounterVec
)

func init() {
	var err error

	requestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_endpoint_milliseconds",
		Help:    "Time taken to process request",
		Buckets: []float64{20, 50, 100, 500},
	}, []string{"endpoint", "method"})

	err = prometheus.Register(requestsDuration)
	if err != nil {
		log.WithError(err).
			WithField("metric", "request_duration_endpoint_milliseconds").
			Error("unable to register prometheus metric")
	}

	requestsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_count_endpoint",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	}, []string{"endpoint", "method"})

	err = prometheus.Register(requestsCounter)
	if err != nil {
		log.WithError(err).
			WithField("metric", "request_count_endpoint").
			Error("unable to register prometheus metric")
	}
}

// Prometheus возвращает Gin middleware для регистрации метрик Prometheus.
func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// detect execution time
		endpoint := c.FullPath()
		method := c.Request.Method
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			requestsDuration.WithLabelValues(endpoint, method).Observe(v * msMultiplier)
		}))

		defer timer.ObserveDuration()

		requestsCounter.WithLabelValues(endpoint, method).Inc()

		c.Next()
	}
}
