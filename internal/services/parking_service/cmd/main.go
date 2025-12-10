package main

import (
	"context"
	"k071123/internal/services/parking_service/core"
	"k071123/internal/services/parking_service/delivery/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := di.Ctx.Services().Logger()

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
		log.Info("http server started...")
	}()
	go func() {
		grpcServer.Start()
		defer wg.Done()
		log.Info("grpc server started...")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Infoln("Shutting down parking service...")
	cancel()
	time.Sleep(time.Second)

	wg.Wait()
	log.Infoln("Parking service stopped gracefully")
}
