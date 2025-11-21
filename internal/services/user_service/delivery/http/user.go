package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/user_service/domain/cases/user"
	_ "k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	"k071123/internal/utils/middleware"
)

type UserHandler struct {
	userUseCase *user.UserUseCase
	mwr         *middleware.Middleware
}

func NewUserHandler(useCase *user.UserUseCase, mwr *middleware.Middleware) *UserHandler {
	return &UserHandler{
		userUseCase: useCase,
		mwr:         mwr,
	}
}

func RegisterUserRoutes(router fiber.Router, mw *middleware.Middleware, uh *UserHandler) {
	router.Post("", mw.AuthMiddleware([]permissions.Permission{
		permissions.CreateUser,
	}), uh.CreateUserHandler)
	router.Put("", mw.AuthMiddleware([]permissions.Permission{}), uh.UpdateProfileHandler)
}

// CreateUserHandler
// @Summary      Создать пользователя
// @Description  Создать пользователя
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.CreateUserReq true "Данные для создания пользователя."
// @Success      200 {object} props.CreateUserResp "Успешный ответ. Модель пользователя."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUserHandler(ctx *fiber.Ctx) error {
	var args props.CreateUserReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.userUseCase.CreateUser(&args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}

// UpdateProfileHandler
// @Description  Обновление профиля пользователя
// @Summary      Обновление профиля пользователя
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param request body props.UpdateProfileReq true "Данные для обновления пользователя."
// @Success      200 {object} props.UpdateProfileResp "Успешный ответ. Модель пользователя."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/users [put]
func (h *UserHandler) UpdateProfileHandler(ctx *fiber.Ctx) error {
	userUUID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	var args props.UpdateProfileReq
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}
	args.UserUUID = userUUID

	resp, err := h.userUseCase.UpdateProfile(args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}
