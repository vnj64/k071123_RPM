package http

import (
	"github.com/gofiber/fiber/v2"
	_ "k071123/internal/services/notification_service/domain/props"
	"k071123/internal/services/user_service/domain/cases/auth"
	"k071123/internal/services/user_service/domain/props"
	"k071123/internal/utils/errs"
)

type AuthHandler struct {
	useCase *auth.AuthUseCase
}

func NewAuthHandler(useCase *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase}
}

func RegisterAuthRoutes(router fiber.Router, uh *AuthHandler) {
	router.Post("/send_code", uh.SendCodeHandler)
	router.Post("/confirm_code", uh.ConfirmCodeHandler)
}

// SendCodeHandler
// @Summary      Отправить код подтверждения
// @Description  Отправить код подтверждения
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.SendCodeReq true "Почта пользователя."
// @Success      200 {object} props.SendCodeResponse "Статус отправки Email."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/auth/send_code [post]
func (h *AuthHandler) SendCodeHandler(ctx *fiber.Ctx) error {
	var args props.SendCodeReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.SendCode(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}

// ConfirmCodeHandler
// @Summary      Подтвердить отправленный код
// @Description  Подтвердить отправленный код
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.ConfirmCodeReq true "Код потверждения и почта."
// @Success      200 {object} props.ConfirmCodeResp "Успешная авторизация."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/auth/confirm_code [post]
func (h *AuthHandler) ConfirmCodeHandler(ctx *fiber.Ctx) error {
	var args props.ConfirmCodeReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.ConfirmCode(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}
