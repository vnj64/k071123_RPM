package http

import (
	"github.com/gofiber/fiber/v2"
	"k071123/internal/services/notification_service/domain/cases/sender"
	"k071123/internal/services/notification_service/domain/props"
	_ "k071123/internal/services/user_service/domain/props"
	"k071123/internal/utils/errs"
)

type EmailSenderHandler struct {
	useCase *sender.EmailSenderUseCase
}

func NewEmailSenderHandler(useCase *sender.EmailSenderUseCase) *EmailSenderHandler {
	return &EmailSenderHandler{
		useCase: useCase,
	}
}

func RegisterEmailSenderRoutes(router fiber.Router, uh *EmailSenderHandler) {
	//router.Post("/email", uh.SendEmailHandler)
}

// SendEmailHandler
// @Summary      Отправить Email
// @Description  Отправить Email
// @Tags         EmailSender
// @Accept       json
// @Produce      json
// @Param request body props.SendEmailRequest true "Данные для отправки email."
// @Success      200 {object} props.SendEmailResp "Статус отправленного email."
// @Failure      400 {object} errs.Error "Bad Request"
// @Failure      404 {object} errs.Error "Profile not found"
// @Failure      500 {object} errs.Error "Internal Server Error"
// @Router       /api/v1/sender/email [post]
func (h *EmailSenderHandler) SendEmailHandler(ctx *fiber.Ctx) error {
	var args props.SendEmailRequest
	if err := ctx.BodyParser(&args); err != nil {
		return errs.SendError(ctx, err)
	}

	resp, err := h.useCase.SendEmail(&args)
	if err != nil {
		return errs.SendError(ctx, err)
	}

	return errs.SendSuccess(ctx, fiber.StatusOK, resp)
}
