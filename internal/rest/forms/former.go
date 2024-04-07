package forms

import (
	"github.com/Fruitfulfriends-REST-API-server/pkg/rest/response"
	"github.com/gin-gonic/gin"
)

type Former interface {
	ParseAndValidate(c *gin.Context) (Former, response.Error)
	ConvertToMap() map[string]interface{}
}
