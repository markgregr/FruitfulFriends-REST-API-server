package rest

import (
	"context"
	"fmt"
	"github.com/Fruitfulfriends-REST-API-server/internal/config"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/handlers"
	"time"

	ginlogrus "github.com/Toorop/gin-logrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	checkingInterval = time.Minute
)

type Worker struct {
	cfgRest     *config.HTTPServer
	logger      *log.Logger
	apiHandlers []handlers.APIHandler
}

func NewWorker(
	cfgRest *config.HTTPServer,
	logger *log.Logger,
	apiHandlers []handlers.APIHandler,
) *Worker {
	w := &Worker{
		cfgRest:     cfgRest,
		logger:      logger,
		apiHandlers: apiHandlers,
	}

	return w
}

func (w *Worker) Start(ctx context.Context) error {
	const op = "rest.Worker.Start"
	log := w.logger.WithField("method", op)

	for {

		if err := w.run(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Info("starting rest worker")
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", op, ctx.Err())
		case <-time.After(checkingInterval):
		}
	}
}

func (w *Worker) run() error {
	const op = "rest.Worker.run"
	log := w.logger.WithField("method", op)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(ginlogrus.Logger(w.logger), gin.Recovery())
	router.Use(gin.Recovery())

	if w.cfgRest.AllowOrigin != "" {
		log.WithField("allow_origin", w.cfgRest.AllowOrigin).Info("setting up CORS")
		router.Use(CORS(w.cfgRest.AllowOrigin))
	}

	for _, h := range w.apiHandlers {
		log.WithField("handler", h).Info("enriching routes")
		h.EnrichRoutes(router)
	}

	w.addRouters(router)
	log.WithField("host", w.cfgRest.Host).WithField("port", w.cfgRest.Port).Info("running rest worker")

	return router.Run(fmt.Sprintf("%s:%d", w.cfgRest.Host, w.cfgRest.Port))
}

func CORS(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Total-Count, X-Offset, X-Limit, X-Next-Page")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
