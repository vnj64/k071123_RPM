package main

import (
	"k071123/internal/services/user_service/core"
	"k071123/internal/services/user_service/delivery/http"
	"k071123/internal/services/user_service/docs"
	"k071123/internal/services/user_service/services/config"
	"sync"
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
	server := core.NewHttpServer()
	cfg := config.Make()

	docs.SwaggerInfo.Host = "127.0.0.1:" + cfg.HttpPort()

	handlers := &http.Handlers{
		UserHandler: di.UserHandler,
		AuthHandler: di.AuthHandler,
	}

	http.SetupRoutes(server.App(), handlers)
	wg.Add(1)
	go func() {
		server.Start()
		defer wg.Done()
	}()
	wg.Wait()

}
