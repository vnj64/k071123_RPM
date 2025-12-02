package payment

import (
	"context"
	"github.com/google/uuid"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/models/payment_statuses"
	"k071123/internal/services/order_service/domain/props"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	"k071123/internal/utils/errs"
	"strconv"
)

type PaymentUseCase struct {
	ctx domain.Context
	pc  proto.ParkingClient
}

func NewPaymentUseCase(ctx domain.Context, pc proto.ParkingClient) *PaymentUseCase {
	return &PaymentUseCase{
		ctx: ctx,
		pc:  pc,
	}
}

func (uc *PaymentUseCase) MakePayment(args props.CreatePaymentReq) (resp props.CreatePaymentResp, err error) {
	if err := args.Validate(); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, err.Error())
	}

	feeInt, err := strconv.Atoi(uc.ctx.Services().Config().PlatformFee())
	feeFloat := float64(feeInt)

	payment := &models.Payment{
		UUID:          uuid.New(),
		SessionUUID:   args.SessionUUID,
		PaymentMethod: args.PaymentMethod,
		Status:        payment_statuses.Succeeded,
		TransactionId: uuid.New().String(),
		Amount:        args.Amount,
		PlatformFee:   &feeFloat,
		Description:   args.Description,
	}
	if err := uc.ctx.Connection().Payment().Insert(payment); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to create payment [database error]")
	}

	cards, err := uc.ctx.Connection().Card().WhereFilter(uc.ctx.Connection().Card().Filter().SetUserUUIDs([]string{args.UserUUID.String()}))
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to find user bank card [database error]")
	}
	if len(cards) == 0 {
		return resp, errs.NewErrorWithDetails(errs.ErrNotFound, "no user cards connected")
	}
	card := cards[0]
	if err := uc.ctx.Services().Billing().AutoPay(*card.Token, args.Amount); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to make payment [billing error]")
	}
	resp.Payment = payment
	return resp, nil
}

func (uc *PaymentUseCase) UpdatePaymentSessionStatus(paymentID string, status payment_statuses.PaymentStatuses) error {
	payment, err := uc.ctx.Connection().Payment().GetByTransactionID(paymentID)
	if err != nil {
		return errs.NewErrorWithDetails(errs.ErrNotFound, "payment not found")
	}

	updates := uc.ctx.Connection().Payment().Updates()
	updates.SetStatus(status)
	if updates.HaveUpdates() {
		err = uc.ctx.Connection().Payment().Update(payment.UUID, updates)
		if err != nil {
			return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to update payment status")
		}
	}

	_, err = uc.pc.UpdateSessionPaid(context.Background(), &proto.UpdateSessionPaidRequest{
		SessionUuid: payment.SessionUUID.String(),
	})
	if err != nil {
		return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to update payment status")
	}

	return nil
}

func (uc *PaymentUseCase) FindPayment(args props.FindPayments) ([]models.Payment, error) {
	filter := uc.ctx.Connection().Payment().Filter()
	if len(args.Statuses) > 0 {
		filter.SetStatuses(args.Statuses)
	}

	if len(args.SessionUUIDs) > 0 {
		filter.SetSessionUUIDs(args.SessionUUIDs)
	}

	payments, err := uc.ctx.Connection().Payment().WhereFilter(filter, nil)
	if err != nil {
		return nil, errs.NewErrorWithDetails(errs.ErrInternalServerError, "cant find payments")
	}
	return payments, nil
}

func ptrFloat64(f *float64) float64 {
	v := f
	return *v
}
