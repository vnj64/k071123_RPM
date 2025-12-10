package http

import (
	"github.com/gofiber/fiber/v2"

	"k071123/internal/services/parking_service/domain/cases/tariff"
	"k071123/internal/services/parking_service/domain/props"

	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
)

type TariffHandler struct {
	useCase *tariff.TariffUseCase
	mw      *middleware.Middleware
}

func NewTariffHandler(useCase *tariff.TariffUseCase, mw *middleware.Middleware) *TariffHandler {
	return &TariffHandler{
		useCase: useCase,
		mw:      mw,
	}
}

func RegisterTariffRoutes(router fiber.Router, ph *TariffHandler, mw *middleware.Middleware) {
	router.Post("/create", mw.AuthMiddleware([]permissions.Permission{
		permissions.CreateUser,
	}), ph.CreateTariffHandler)
}

// CreateTariffHandler
// @Summary      Создать тариф
// @Description  Создать тариф
// @Tags         Tariff
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.CreateTariffReq true "Данные для создания тарифа"
// @Success      200 {object} props.CreateTariffResp "Модель созданной тарифа."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/tariff/create [post]
func (h *TariffHandler) CreateTariffHandler(ctx *fiber.Ctx) error {
	var args props.CreateTariffReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.CreateTariff(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}
