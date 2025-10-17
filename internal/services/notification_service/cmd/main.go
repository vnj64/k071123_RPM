package main

import (
	"k071123/internal/services/notification_service/core"
	"k071123/internal/services/notification_service/delivery/http"
	"k071123/internal/services/notification_service/docs"
	"k071123/internal/services/notification_service/services/config"
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

	cfg := config.Make()
	docs.SwaggerInfo.Host = "127.0.0.1:" + cfg.HttpPort()

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
