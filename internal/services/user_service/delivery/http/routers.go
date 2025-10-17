package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes регистрирует все маршруты, используя агрегированный набор хэндлеров
// @title User Service API
// @version 1.0
// @description API для работы с пользователями
// @host localhost:7800
// @BasePath /api/v1
// @schemes http
func SetupRoutes(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	docs := app.Group("/docs")
	docs.Get("/swagger/*", swagger.HandlerDefault)

	RegisterUserRoutes(api.Group("/users"), h.UserHandler.mwr, h.UserHandler)
	RegisterAuthRoutes(api.Group("/auth"), h.AuthHandler)
}
