package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil) // disable trusted proxies warning
	router.Use(gin.Recovery())        // recovery from panic()
	router.GET("/health", func(c *gin.Context) {
		// TODO : add a check for RAFT connection?
		c.JSON(200, gin.H{"status": "ok"})
	})
	return router
}
