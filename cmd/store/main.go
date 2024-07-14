package main

import (
	"codesignal/internal/api"
	"codesignal/internal/server"
	"codesignal/internal/services"
	"codesignal/internal/store"
	"log"
	"os"
	"os/signal"
)

func main() {

	log.Println("Starting Raft...")
	kvstore := store.SetupStore()
	keyValueStore, err := services.NewRaftKeyValueStore(kvstore)
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
	router = server.SetupApiRoutes(router, api.NewKeyValueStoreApi(keyValueStore))
	router = server.SetupRaftRoutes(router, kvstore)

	err = router.Run(":" + PORT)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
}
