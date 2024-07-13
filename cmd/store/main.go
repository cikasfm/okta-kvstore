package main

import (
	"codesignal/internal/api"
	"codesignal/internal/raft"
	"codesignal/internal/server"
	"codesignal/internal/services"
	"log"
	"os"
)

func main() {

	log.Println("Starting Raft...")
	ra := raft.SetupRaft()
	keyValueStore, err := services.NewRaftKeyValueStore(ra)
	if err != nil {
		panic(err)
	}
	log.Println("Starting Raft -> OK")

	//gin.SetMode(gin.ReleaseMode)
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	log.Println("Starting server on port", PORT)

	// create router
	router := server.SetupRouter()
	router = server.SetupRoutes(router, api.NewKeyValueStoreApi(keyValueStore))
	router.POST("/join", server.JoinHandler(ra))

	err = router.Run(":" + PORT)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
}
