package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes регистрирует все маршруты, используя агрегированный набор хэндлеров
// @title Notification Service API
// @version 1.0
// @description API для рассылки уведомлений
// @host localhost:7801
// @BasePath /api/v1
// @schemes http
func SetupRoutes(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	docs := app.Group("/docs")
	docs.Get("/swagger/*", swagger.HandlerDefault)

	sender := api.Group("/sender")

	// Email Sender
	RegisterEmailSenderRoutes(sender, h.EmailSenderHandler)
}
