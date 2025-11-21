package main

import (
	"k071123/internal/services/notification_service/core"
	"k071123/internal/services/notification_service/delivery/http"
	"log"
	"sync"
)

// @title Notification Service API
// @version 2.0
// @description Сервис для рассылки уведомлений
// @host localhost:7801
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
		EmailSenderHandler: di.EmailSenderHandler,
	}

	http.SetupRoutes(server.App(), handlers)
	wg.Add(2)
	go func() {
		server.Start()
		log.Print("Http server starting...")
		defer wg.Done()
	}()
	go func() {
		grpcServer.Start()
		log.Print("GRPC server starting...")
		defer wg.Done()
	}()
	wg.Wait()

}
