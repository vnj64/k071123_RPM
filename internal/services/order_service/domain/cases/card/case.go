package card

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"k071123/internal/services/notification_service/contracts/pkg/proto"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/props"
	"k071123/internal/utils/errs"
	"math/rand"
	"strconv"
)

type CardUseCase struct {
	ctx domain.Context
	nc  proto.NotificationClient
}

func NewCardUseCase(ctx domain.Context, nc proto.NotificationClient) *CardUseCase {
	return &CardUseCase{
		ctx: ctx,
		nc:  nc,
	}
}

func (uc *CardUseCase) GetCardByUUID(args props.GetCard) (*models.Card, error) {
	if err := args.Validate(); err != nil {
		return nil, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, fmt.Sprintf("validation error: %v", err))
	}
	result, err := uc.ctx.Connection().Card().GetByUUID(args.UUID)
	if err != nil {
		return nil, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on insert to database")
	}
	return result, nil
}

func (uc *CardUseCase) SaveCard(args props.SaveCardReq) (resp props.SaveCardResp, err error) {
	if err := args.Validate(); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, fmt.Sprintf("validation error: %v", err))
	}

	card := &models.Card{
		UUID:          uuid.New(),
		Last4Digits:   getLast4(args.CardNumber),
		PaymentSystem: args.PaymentSystem,
		UserUUID:      args.UserUUID,
		IsActive:      false,
	}
	if args.IsPreferred {
		card.IsPreferred = args.IsPreferred
		if err := uc.ctx.Connection().Card().ChangePreferredCard(args.UserUUID, *card); err != nil {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on update card")
		}
	} else {
		if err := uc.ctx.Connection().Card().Insert(card); err != nil {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on insert card to database")
		}
	}

	otp := generateRandomToken()
	if err := uc.ctx.Connection().VerifyToken().Insert(&models.VerifyTokens{
		UUID:     uuid.New(),
		UserUUID: uuid.MustParse(args.UserUUID),
		Token:    otp,
		Used:     false,
	}); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on save otp")
	}

	notificationResp, err := uc.nc.SendEmail(context.Background(), &proto.SendEmailReq{
		Subject: "Card OTP Verifying",
		Data:    generateRandomToken(),
		To:      []string{args.Email},
	})
	if notificationResp != nil {
		if notificationResp.Response != "success" {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to send email")
		}
	}
	resp.Message = fmt.Sprintf("We have sent you OTP to email %s\n. Please verify that your card.", args.Email)

	return resp, nil
}

func (uc *CardUseCase) VerifyCard(args props.VerifyCardReq) (resp props.VerifyCardResp, err error) {
	otp, err := uc.ctx.Connection().VerifyToken().GetLastByUserUUID(args.UserUUID)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on verify otp")
	}
	if otp.Token != args.OTP {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "wrong otp code")
	}
	//if err := uc.ctx.Connection().Card().
	return resp, nil
}

func getLast4(cardNumber string) string {
	return cardNumber[:4]
}

func generateRandomToken() string {
	return strconv.Itoa(rand.Int())[:4]
}
