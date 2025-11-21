package core

import (
	"k071123/internal/services/order_service/delivery/http"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/cases/card"
	"k071123/internal/services/order_service/services/config"
	"k071123/internal/services/user_service/core"
	"k071123/internal/utils/middleware"
	"log"
)

type Di struct {
	Ctx         domain.Context
	CardHandler *http.CardHandler
}

func NewDi() *Di {
	cfg := config.Make()
	ctx := InitCtx()

	nc, err := core.MakeNotificationServiceClient()
	if err != nil {
		log.Fatalf("create parking service client failed: %v", err)
	}
	mw := middleware.NewMiddleware(cfg.PublicPemPath())
	var (
		cardUseCase = card.NewCardUseCase(ctx, nc)
		cardHandler = http.NewCardHandler(cardUseCase, mw)
	)

	return &Di{
		Ctx:         ctx,
		CardHandler: cardHandler,
	}
}
