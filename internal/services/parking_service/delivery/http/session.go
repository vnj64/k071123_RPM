package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain/cases/car"
	"k071123/internal/services/parking_service/domain/cases/session"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
	"k071123/tools/logger"
)

type SessionHandler struct {
	useCase    *session.SessionUseCase
	carUseCase *car.CarUseCase
	mw         *middleware.Middleware
	oc         proto.OrderClient
	log        *logger.Logger
}

func NewSessionHandler(
	useCase *session.SessionUseCase,
	carUc *car.CarUseCase,
	mw *middleware.Middleware,
	oc proto.OrderClient,
	log *logger.Logger,
) *SessionHandler {
	return &SessionHandler{
		useCase:    useCase,
		carUseCase: carUc,
		mw:         mw,
		oc:         oc,
		log:        log,
	}
}

func RegisterSessionRoutes(router fiber.Router, sh *SessionHandler, mw *middleware.Middleware) {
	router.Post("/start", mw.AuthMiddleware([]permissions.Permission{
		permissions.StartSession,
	}), sh.StartSessionHandler)
	router.Post("/stop", mw.AuthMiddleware([]permissions.Permission{
		permissions.StopSession,
	}), sh.FinishSessionHandler)
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
	log := h.log.WithField("Handler", "StartSession")
	if err := ctx.BodyParser(&args); err != nil {
		log.Errorf("unable to parse StartSession request body: %v", err)
		return errs.SendError(ctx, err)
	}

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		log.Errorf("unable to get userUUID from context: %v", err)
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

	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		log.Infof("user: %s", userUUID)
		return errs.SendError(ctx, err)
	}
	log.Infof("car number: %s", args.CarNumber)
	args.UserUUID = userUUID
	resp, err := h.useCase.Finish(args, h.oc)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusCreated, resp)
}
