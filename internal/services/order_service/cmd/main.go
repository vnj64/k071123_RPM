package main

import (
	"context"
	"k071123/internal/services/order_service/core"
	"k071123/internal/services/order_service/delivery/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	log := di.Ctx.Services().Logger().WithField("HttpServer", "Order")
	server := core.NewHttpServer()
	grpcServer := core.NewGrpcServer()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlers := &http.Handlers{
		CardHandler: di.CardHandler,
	}

	http.SetupRoutes(server.App(), handlers)
	wg.Add(2)
	go func() {
		server.Start()
		defer wg.Done()
	}()
	log.Info("Http Order server is starting...")
	go func() {
		grpcServer.Start()
		defer wg.Done()
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Infoln("Shutting down order service...")
	cancel()
	time.Sleep(time.Second)

	wg.Wait()
	log.Infoln("Order service stopped gracefully")
}
