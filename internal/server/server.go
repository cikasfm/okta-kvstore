package server

import (
	"codesignal/internal/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil) // disable trusted proxies warning
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, "OK")
	})
	return router
}

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

func setupRoutes(router *gin.Engine, kvStoreRestApi api.IKeyValueStoreApi) *gin.Engine {
	// setup routes
	router = getKey(router, kvStoreRestApi.GetByKey)
	router = postKey(router, kvStoreRestApi.SetValue)
	router = deleteKey(router, kvStoreRestApi.DeleteValue)
	return router
}

func StartServer(kvStoreRestApi api.IKeyValueStoreApi, addr ...string) {
	// gin.SetMode(gin.ReleaseMode)

	fmt.Println("Starting server on", addr)

	// create router
	router := setupRouter()
	router = setupRoutes(router, kvStoreRestApi)

	err := router.Run(addr...)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
