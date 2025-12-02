package main

import (
	"k071123/internal/services/order_service/core"
	"k071123/internal/services/order_service/delivery/http"
	"sync"
)

// @title Order Service API
// @version 2.0
// @description API для работы с доменом Order
// @host localhost:7803
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Enter your Bearer token in the format: `Bearer {token}`"
func main() {
	var wg sync.WaitGroup

	di := core.NewDi()
	server := core.NewHttpServer()
	grpcServer := core.NewGrpcServer()

	handlers := &http.Handlers{
		CardHandler: di.CardHandler,
	}

	http.SetupRoutes(server.App(), handlers)
	wg.Add(2)
	go func() {
		server.Start()
		defer wg.Done()
	}()
	go func() {
		grpcServer.Start()
		defer wg.Done()
	}()
	wg.Wait()
}
