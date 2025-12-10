package sender

import (
	"context"
	"k071123/internal/services/notification_service/domain"
	"k071123/internal/services/notification_service/domain/models"
	"k071123/internal/services/notification_service/domain/props"
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
	log := uc.ctx.Services().Logger().WithField("EmailSenderUseCase", "SendEmail")
	email := models.Email{
		Subject: args.Subject,
		From:    uc.ctx.Services().Config().SmtpFrom(),
		Data:    args.Data,
	}
	if err := uc.ctx.Services().Amqp().SendEmail(context.Background(), email); err != nil {
		if err := uc.ctx.Services().Smtp().Send(&email); err != nil {
			resp.Status = "success"
			log.WithError(err).Infof("send email err %s", err.Error())
			return resp, err
		}
		resp.Status = "failed"
		log.WithError(err).Infof("send email err %s", err.Error())
		return resp, err
	}
	resp.Status = "queued"

	return resp, nil
}
