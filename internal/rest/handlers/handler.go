package handlers

import "github.com/gin-gonic/gin"

type APIHandler interface {
	EnrichRoutes(router *gin.Engine)
}
