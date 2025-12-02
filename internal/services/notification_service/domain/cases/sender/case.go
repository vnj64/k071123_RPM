package sender

import (
	"k071123/internal/services/notification_service/domain"
	"k071123/internal/services/notification_service/domain/models"
	"k071123/internal/services/notification_service/domain/props"
	"k071123/internal/utils/errs"
)

type EmailSenderUseCase struct {
	ctx domain.Context
}

func NewEmailSenderUseCase(ctx domain.Context) *EmailSenderUseCase {
	return &EmailSenderUseCase{
		ctx: ctx,
	}
}

func (uc *EmailSenderUseCase) SendEmail(args *props.SendEmailRequest) (resp props.SendEmailResp, err error) {
	if err := uc.ctx.Services().Smtp().Send(&models.Email{
		Subject: args.Subject,
		From:    uc.ctx.Services().Config().SmtpFrom(),
		Data:    args.Data,
	}); err != nil {
		resp.Status = "failed"
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, err.Error())
	}
	// TODO: rabbit subscriber
	resp.Status = "success"

	return resp, nil
}
