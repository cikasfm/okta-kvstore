package main

import (
	"codesignal/internal/api"
	"codesignal/internal/server"
	"codesignal/internal/services"
)

func main() {
	keyValueStore, err := services.NewKeyValueStore()
	if err != nil {
		panic(err)
	}
	server.StartServer(api.NewKeyValueStoreApi(keyValueStore), ":8080")
}
