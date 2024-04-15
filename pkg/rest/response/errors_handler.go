package response

import (
	"github.com/gin-gonic/gin"
)

func HandleError(err Error, c *gin.Context) {
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
