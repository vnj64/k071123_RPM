package payment

import (
	"context"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/models/payment_statuses"
	"k071123/internal/services/order_service/domain/props"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	"k071123/internal/utils/errs"
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
