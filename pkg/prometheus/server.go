package prometheus

import (
	"github.com/Fruitfulfriends-REST-API-server/pkg/middleware"
	ginlogrus "github.com/Toorop/gin-logrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NewPrometheusServer(log *logrus.Logger, listen, path string) *http.Server {
	const op = "pkg.prometheus.NewPrometheusServer"
	log.WithField("method", op)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.Panic(log), middleware.Prometheus(), ginlogrus.Logger(log))

	r.GET(path, gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:    listen,
		Handler: r,
	}

	return server
}