package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain/cases/car"
	"k071123/internal/services/parking_service/domain/cases/session"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
	"log"
)

type SessionHandler struct {
	useCase    *session.SessionUseCase
	carUseCase *car.CarUseCase
	mw         *middleware.Middleware
	oc         proto.OrderClient
}

func NewSessionHandler(
	useCase *session.SessionUseCase,
	carUc *car.CarUseCase,
	mw *middleware.Middleware,
	oc proto.OrderClient,
) *SessionHandler {
	return &SessionHandler{
		useCase:    useCase,
		carUseCase: carUc,
		mw:         mw,
		oc:         oc,
	}
}

func RegisterSessionRoutes(router fiber.Router, sh *SessionHandler, mw *middleware.Middleware) {
	router.Post("/start", mw.AuthMiddleware([]permissions.Permission{
		permissions.StartSession,
	}), sh.StartSessionHandler)
	router.Post("/stop", mw.AuthMiddleware([]permissions.Permission{
		permissions.StopSession,
	}))
	router.Post("/car", mw.AuthMiddleware([]permissions.Permission{}), sh.CreateCarHandler)
}

// StartSessionHandler
// @Summary      Начать парковочную сессию
// @Description  Начать парковочную сессию
// @Tags         Session
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.StartSessionReq true "Данные для старта сесси"
// @Success      200 {object} props.StartSessionResp "Статус сессии."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/session/start [post]
func (h *SessionHandler) StartSessionHandler(ctx *fiber.Ctx) error {
	var args props.StartSessionReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	args.UserUUID = userUUID
	resp, err := h.useCase.Start(args, h.oc)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}

// FinishSessionHandler
// @Summary      Закончить парковочную сессию
// @Description  Закончить парковочную сессию
// @Tags         Session
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.FinishSessionRequest true "Данные для стопа сессии"
// @Success      200 {object} props.FinishSessionResp "Статус сессии."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/session/stop [post]
func (h *SessionHandler) FinishSessionHandler(ctx *fiber.Ctx) error {
	var args props.FinishSessionRequest
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.Finish(args, h.oc)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}

// CreateCarHandler
// @Summary      Закончить парковочную сессию
// @Description  Закончить парковочную сессию
// @Tags         Session
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.CreateCarReq true "Данные для стопа сессии"
// @Success      200 {object} props.CreateCarResp "Статус сессии."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/session/car [post]
func (h *SessionHandler) CreateCarHandler(ctx *fiber.Ctx) error {
	var args props.CreateCarReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return errs.SendError(ctx, err)
	}
	args.UserUUID = userUUID
	log.Printf("user uuid from ctx %s", userUUID)

	resp, err := h.carUseCase.CreateCar(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}
