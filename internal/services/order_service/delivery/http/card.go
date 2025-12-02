package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/order_service/domain/cases/card"
	"k071123/internal/services/order_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
)

type CardHandler struct {
	useCase *card.CardUseCase
	mw      *middleware.Middleware
}

func NewCardHandler(useCase *card.CardUseCase, mw *middleware.Middleware) *CardHandler {
	return &CardHandler{
		useCase: useCase,
		mw:      mw,
	}
}

func RegisterCardRoutes(router fiber.Router, ph *CardHandler, mw *middleware.Middleware) {
	router.Post("/save", mw.AuthMiddleware([]permissions.Permission{
		permissions.SaveCard,
	}), ph.SaveCardHandler)
	router.Post("/verify", mw.AuthMiddleware([]permissions.Permission{
		permissions.VerifyCard,
	}), ph.VerifyCardHandler)
}

// SaveCardHandler
// @Summary      Сохранить банковскую карту
// @Description  Сохранить банковскую карту
// @Tags         Card
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.SaveCardReq true "Данные для создания парковки"
// @Success      200 {object} props.SaveCardResp "Success message"
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/card/save [post]
func (h *CardHandler) SaveCardHandler(ctx *fiber.Ctx) error {
	var args props.SaveCardReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	args.UserUUID = userUUID
	resp, err := h.useCase.SaveCard(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}

// VerifyCardHandler
// @Summary      Подтвердить банковскую карту
// @Description  Подтвердить банковскую карту
// @Tags         Card
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.VerifyCardReq true "Данные для подтверждения карты"
// @Success      200 {object} props.VerifyCardResp "Success message"
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/card/verify [post]
func (h *CardHandler) VerifyCardHandler(ctx *fiber.Ctx) error {
	var args props.VerifyCardReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return errs.SendError(ctx, err)
	}
	args.UserUUID = userUUID

	resp, err := h.useCase.VerifyCard(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}
