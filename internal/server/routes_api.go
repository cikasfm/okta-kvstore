package server

import (
	"codesignal/internal/api"
	"github.com/gin-gonic/gin"
)

// GET /key/:key
func getKey(router *gin.Engine, handler gin.HandlerFunc) *gin.Engine {
	router.GET("/key/:key", handler)
	return router
}

// POST /key
func postKey(router *gin.Engine, handler gin.HandlerFunc) *gin.Engine {
	router.POST("/key", handler)
	return router
}

// DELETE /key/:key
func deleteKey(router *gin.Engine, handler gin.HandlerFunc) *gin.Engine {
	router.DELETE("/key/:key", handler)
	return router
}

func SetupRoutes(router *gin.Engine, kvStoreRestApi api.IKeyValueStoreApi) *gin.Engine {
	// setup routes
	router = getKey(router, kvStoreRestApi.GetByKey)
	router = postKey(router, kvStoreRestApi.SetValue)
	router = deleteKey(router, kvStoreRestApi.DeleteValue)
	return router
}
