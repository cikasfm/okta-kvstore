package main

import "codesignal/cmd/store/services"

func main() {
	service, err := services.NewKeyValueStore()
	if err != nil {
		panic(err)
	}
	startServer(NewHttpHandler(service))
}
