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

func (uc *CardUseCase) SaveCard(args props.SaveCardReq) (resp props.SaveCardResp, err error) {
	log := uc.ctx.Services().Logger().WithField("CardUseCase", "SaveCard")
	if err := args.Validate(); err != nil {
		log.Errorf("validation error on card: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, fmt.Sprintf("validation error: %v", err))
	}

	card := &models.Card{
		UUID:          uuid.New(),
		Last4Digits:   getLast4(args.CardNumber),
		PaymentSystem: args.PaymentSystem,
		UserUUID:      args.UserUUID,
		IsActive:      false,
	}
	token := uc.ctx.Services().Billing().GeneratePayToken(card.Last4Digits, args.CVC, args.Date)
	log.Infof("pay token for user %s is successfully generated", card.UserUUID)
	card.Token = &token
	if args.IsPreferred {
		card.IsPreferred = args.IsPreferred
		if err := uc.ctx.Connection().Card().ChangePreferredCard(args.UserUUID, *card); err != nil {
			log.Errorf("unable to change preferred card [db]")
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on update card")
		}
	} else {
		if err := uc.ctx.Connection().Card().Insert(card); err != nil {
			log.Errorf("unable to insert card [db]")
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on insert card to database")
		}
	}

	otp := generateRandomToken()
	if err := uc.ctx.Connection().VerifyToken().Insert(&models.VerifyTokens{
		UUID:     uuid.New(),
		UserUUID: uuid.MustParse(args.UserUUID),
		OTP:      otp,
		Used:     false,
	}); err != nil {
		log.Errorf("unable to save otp for card: %v [db]", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on save otp")
	}

	notificationResp, err := uc.nc.SendEmail(context.Background(), &proto.SendEmailReq{
		Subject: "Card OTP Verifying",
		Data:    otp,
		To:      []string{args.Email},
	})
	if notificationResp != nil {
		if notificationResp.Response != "success" && notificationResp.Response != "queued" {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to send email")
		}
	}
	resp.Message = fmt.Sprintf("We have sent you OTP to email %s\n. Please verify that your card.", args.Email)
	log.Infof("user payment card is save. user: %s", card.UserUUID)

	return resp, nil
}

func (uc *CardUseCase) VerifyCard(args props.VerifyCardReq) (resp props.VerifyCardResp, err error) {
	log := uc.ctx.Services().Logger().WithField("CardUseCase", "VerifyCard")
	otp, err := uc.ctx.Connection().VerifyToken().GetLastByUserUUID(args.UserUUID)
	if err != nil {
		log.Errorf("unable to get verify token: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on verify otp")
	}
	if otp.OTP != args.OTP {
		log.Errorf("otp code is invalid")
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "wrong otp code")
	}

	otp.Used = true
	if err := uc.ctx.Connection().VerifyToken().Save(otp); err != nil {
		log.Errorf("unable to save otp for card: %v [db]", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on update otp")
	}

	filter := uc.ctx.Connection().Card().Filter().SetUserUUIDs([]string{args.UserUUID})
	cards, err := uc.ctx.Connection().Card().WhereFilter(filter)
	if err != nil {
		log.Errorf("unable to find card by uuid: %v", args.UserUUID)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on find cards")
	}
	if len(cards) == 0 {
		log.Errorf("card not found by uuid: %v", args.UserUUID)
		return resp, errs.NewErrorWithDetails(errs.ErrNotFound, "unable to find card")
	}
	card := cards[0]

	card.IsActive = true
	if err := uc.ctx.Connection().Card().Save(&card); err != nil {
		log.Errorf("unable to save card [db]")
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on update card")
	}
	return resp, nil
}

func (uc *CardUseCase) GetPreferredByUserUUID(args props.GetPreferredByUserUUIDReq) (resp props.GetPreferredByUserUUIDResp, err error) {
	log := uc.ctx.Services().Logger().WithField("CardUseCase", "GetPreferredByUserUUID")
	filter := uc.ctx.Connection().Card().Filter().SetUserUUIDs([]string{args.UserUUID}).SetIsPreferred(true)
	cards, err := uc.ctx.Connection().Card().WhereFilter(filter)
	if err != nil {
		log.Errorf("unable to find card by uuid: %v", args.UserUUID)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "error on find cards")
	}
	if len(cards) == 0 {
		log.Errorf("card not found by uuid: %v", args.UserUUID)
		return resp, errs.NewErrorWithDetails(errs.ErrNotFound, "unable to find card")
	}
	card := cards[0]
	resp.Card = &card

	return resp, nil
}

func getLast4(cardNumber string) string {
	return cardNumber[:4]
}

func generateRandomToken() string {
	return strconv.Itoa(rand.Int())[:4]
}
