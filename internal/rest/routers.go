package rest

import (
	"github.com/gin-gonic/gin"
)

func (w *Worker) addRouters(router *gin.Engine) {
	const op = "rest.Worker.addRouters"
	log := w.logger.WithField("method", op)
	log.Info("adding routers")

	router.NoRoute(func(c *gin.Context) {
		log.Info("%s %s\n", c.Request.Method, c.Request.URL)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}
