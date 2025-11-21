package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/parking_service/domain/cases/unit"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
)

type UnitHandler struct {
	useCase *unit.UnitUseCase
	mw      *middleware.Middleware
}

func NewUnitHandler(useCase *unit.UnitUseCase, mw *middleware.Middleware) *UnitHandler {
	return &UnitHandler{
		useCase: useCase,
		mw:      mw,
	}
}

func RegisterUnitRoutes(router fiber.Router, ph *UnitHandler, mw *middleware.Middleware) {
	router.Post("/create", mw.AuthMiddleware([]permissions.Permission{
		permissions.CreateUnit,
	}), ph.CreateUnitHandler)
}

// CreateUnitHandler
// @Summary      Создать сущность Unit
// @Description  Создать сущность Unit
// @Tags         Unit
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.CreateUniqReq true "Данные для создания паркомата"
// @Success      200 {object} props.CreateUnitResp "Модель созданной паркомата."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/unit/create [post]
func (h *UnitHandler) CreateUnitHandler(ctx *fiber.Ctx) error {
	var args props.CreateUniqReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.Create(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}
