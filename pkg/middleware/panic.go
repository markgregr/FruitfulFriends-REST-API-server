package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
)

func Panic(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if rec := recover(); rec != nil {
				var message string

				switch v := rec.(type) {
				case string:
					message = v
				case error:
					message = v.Error()
				default:
					message = "unknown error"
				}

				body, _ := ioutil.ReadAll(c.Request.Body)

				log.WithFields(logrus.Fields{
					"message": message,
					"request": fmt.Sprintf("%s %s?%s", c.Request.Method, c.Request.URL.String(), c.Request.URL.Query().Encode()),
					"body":    string(body),
					"stack":   string(debug.Stack()),
				}).Error("panic raised")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}(c)

		c.Next()
	}
}
