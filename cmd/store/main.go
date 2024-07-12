package main

import (
	"codesignal/cmd/store/api"
	"codesignal/cmd/store/services"
)

func main() {
	keyValueStore, err := services.NewKeyValueStore()
	if err != nil {
		panic(err)
	}
	startServer(api.NewKeyValueStoreApi(keyValueStore), "localhost:8080")
}
