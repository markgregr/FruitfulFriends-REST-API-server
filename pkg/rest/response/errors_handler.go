package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(err Error, c *gin.Context) {
	if err == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(err.GetHTTPStatus(), ParseError(err))
}

// ParseError determines the error type and creates a map with the error description.
func ParseError(err Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *ValidationError:
		return map[string]interface{}{
			"errors":  e.Errors(),
			"message": e.PublicMessage(),
		}
	default:
		return map[string]interface{}{
			"message": err.PublicMessage(),
		}
	}
}
