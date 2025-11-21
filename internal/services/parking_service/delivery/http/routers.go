package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes регистрирует все маршруты, используя агрегированный набор хэндлеров
// @title Parking Service API
// @version 1.0
// @description API для работы с парковками
// @host localhost:7802
// @BasePath /api/v1
// @schemes http
func SetupRoutes(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	docs := app.Group("/docs")
	docs.Get("/swagger/*", swagger.HandlerDefault)

	parkingGroup := api.Group("/parking")
	tariffGroup := api.Group("/tariff")
	sessionGroup := api.Group("/session")
	unitGroup := api.Group("/unit")

	RegisterParkingRoutes(parkingGroup, h.ParkingHandler, h.ParkingHandler.mw)
	RegisterTariffRoutes(tariffGroup, h.TariffHandler, h.TariffHandler.mw)
	RegisterSessionRoutes(sessionGroup, h.SessionHandler, h.SessionHandler.mw)
	RegisterUnitRoutes(unitGroup, h.UnitHandler, h.UnitHandler.mw)
}
