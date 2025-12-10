package main

import (
	"context"
	"k071123/internal/services/user_service/core"
	"k071123/internal/services/user_service/delivery/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title User Service API
// @version 2.0
// @description API для работы с пользователем
// @host localhost:7800
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Enter your Bearer token in the format: `Bearer {token}`"
func main() {
	var wg sync.WaitGroup

	di := core.NewDi()
	log := di.Ctx.Services().Logger()
	server := core.NewHttpServer()

	grpcServer := core.NewGrpcServer()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlers := &http.Handlers{
		UserHandler: di.UserHandler,
		AuthHandler: di.AuthHandler,
	}

	http.SetupRoutes(server.App(), handlers)
	wg.Add(2)
	go func() {
		server.Start()
		defer wg.Done()
		log.Infoln("User HTTP Server starting...")
	}()
	go func() {
		grpcServer.Start()
		defer wg.Done()
		log.Infoln("User gRPC Server starting...")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Infoln("Shutting down user service...")
	cancel()
	time.Sleep(time.Second)

	wg.Wait()
	log.Infoln("User service stopped gracefully")
}
