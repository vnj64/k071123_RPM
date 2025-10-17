package core

import (
	"k071123/internal/services/user_service/delivery/http"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/cases/auth"
	"k071123/internal/services/user_service/domain/cases/user"
	"k071123/internal/services/user_service/services/config"
	"k071123/internal/utils/middleware"
)

type Di struct {
	Ctx         domain.Context
	UserHandler *http.UserHandler
	AuthHandler *http.AuthHandler
}

func NewDi() *Di {
	cfg := config.Make()
	ctx := InitCtx()
	mw := middleware.NewMiddleware(cfg.PublicPemPath())

	notificationGrpcClient, err := MakeNotificationServiceClient()
	if err != nil {
		panic(err)
	}

	var (
		userUseCase = user.NewUserUseCase(ctx)
		userHandler = http.NewUserHandler(userUseCase, mw)

		authUseCase = auth.NewAuthUseCase(ctx, notificationGrpcClient)
		authHandler = http.NewAuthHandler(authUseCase)
	)
	return &Di{
		Ctx:         ctx,
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
