package core

import (
	"k071123/internal/services/notification_service/delivery/http"
	"k071123/internal/services/notification_service/domain"
	"k071123/internal/services/notification_service/domain/cases/sender"
	"k071123/internal/services/notification_service/services/config"
)

type Di struct {
	Ctx                domain.Context
	EmailSenderHandler *http.EmailSenderHandler
}

func NewDi() *Di {
	_ = config.Make()
	ctx := InitCtx()

	var (
		senderUseCase      = sender.NewEmailSenderUseCase(ctx)
		emailSenderHandler = http.NewEmailSenderHandler(senderUseCase)
	)

	var ()
	return &Di{
		Ctx:                ctx,
		EmailSenderHandler: emailSenderHandler,
	}
}
