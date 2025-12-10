package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/parking_service/domain/cases/parking"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
	"k071123/tools/logger"
)

type ParkingHandler struct {
	useCase *parking.ParkingUseCase
	mw      *middleware.Middleware
	log     *logger.Logger
}

func NewParkingHandler(useCase *parking.ParkingUseCase, mw *middleware.Middleware, log *logger.Logger) *ParkingHandler {
	return &ParkingHandler{
		useCase: useCase,
		mw:      mw,
		log:     log,
	}
}

func RegisterParkingRoutes(router fiber.Router, ph *ParkingHandler, mw *middleware.Middleware) {
	router.Post("/create", mw.AuthMiddleware([]permissions.Permission{
		permissions.CreateParking,
	}), ph.CreateParkingHandler)
}

// CreateParkingHandler
// @Summary      Создать сущность парковки
// @Description  Создать сущность парковки
// @Tags         Parking
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.CreateParkingReq true "Данные для создания парковки"
// @Success      200 {object} props.CreateParkingResp "Модель созданной парковки."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/parking/create [post]
func (h *ParkingHandler) CreateParkingHandler(ctx *fiber.Ctx) error {
	var args props.CreateParkingReq
	log := h.log.WithField("Handler", "CreateParking")
	if err := ctx.BodyParser(&args); err != nil {
		log.Errorf("unable to parse CreateParking request body: %v", err)
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.CreateParking(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}
