package core

import (
	"k071123/internal/services/parking_service/delivery/http"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/cases/car"
	"k071123/internal/services/parking_service/domain/cases/parking"
	"k071123/internal/services/parking_service/domain/cases/session"
	"k071123/internal/services/parking_service/domain/cases/tariff"
	"k071123/internal/services/parking_service/domain/cases/unit"
	"k071123/internal/services/parking_service/services/config"
	"k071123/internal/utils/middleware"
)

type Di struct {
	Ctx            domain.Context
	ParkingHandler *http.ParkingHandler
	TariffHandler  *http.TariffHandler
	SessionHandler *http.SessionHandler
	UnitHandler    *http.UnitHandler
}

func NewDi() *Di {
	cfg := config.Make()
	ctx := InitCtx()

	mw := middleware.NewMiddleware(cfg.PublicPemPath())
	var (
		parkingUseCase = parking.NewParkingUseCase(ctx)
		parkingHandler = http.NewParkingHandler(parkingUseCase, mw)

		tariffUseCase = tariff.NewTariffUseCase(ctx)
		tariffHandler = http.NewTariffHandler(tariffUseCase, mw)

		carUseCase     = car.NewCarUseCase(ctx)
		sessionUseCase = session.NewSessionUseCase(ctx)
		sessionHandler = http.NewSessionHandler(sessionUseCase, carUseCase, mw)

		unitUseCase = unit.NewUnitUseCase(ctx)
		unitHandler = http.NewUnitHandler(unitUseCase, mw)
	)

	return &Di{
		Ctx:            ctx,
		ParkingHandler: parkingHandler,
		TariffHandler:  tariffHandler,
		SessionHandler: sessionHandler,
		UnitHandler:    unitHandler,
	}
}
