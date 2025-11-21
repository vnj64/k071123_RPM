package main

import (
	"k071123/internal/services/parking_service/core"
	"k071123/internal/services/parking_service/delivery/http"
	"log"
	"sync"
)

// @title Parking Service API
// @version 2.0
// @description API для работы с доменом Parking
// @host localhost:7802
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
		ParkingHandler: di.ParkingHandler,
		TariffHandler:  di.TariffHandler,
		SessionHandler: di.SessionHandler,
		UnitHandler:    di.UnitHandler,
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
	log.Printf("grpc server started")
	wg.Wait()
}
