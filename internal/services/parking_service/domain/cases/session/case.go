package session

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/cases/tariff"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
	"time"
)

type SessionUseCase struct {
	ctx domain.Context
	oc  proto.OrderClient
}

func NewSessionUseCase(ctx domain.Context, oc proto.OrderClient) *SessionUseCase {
	return &SessionUseCase{ctx: ctx, oc: oc}
}

func (uc *SessionUseCase) Start(args props.StartSessionReq, oc proto.OrderClient) (resp props.StartSessionResp, err error) {
	// unit parking uuid validation
	cardResp, err := oc.GetPreferredByUserUUID(context.Background(), &proto.GetPreferredCardReq{
		UserUuid: args.UserUUID,
	})
	if err != nil {
		log.Errorf("err: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get user card from order service")
	}
	if cardResp.Card == nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "cannot start session without connected card")
	}

	unit, err := uc.ctx.Connection().UnitRepository().GetByUUID(args.UnitUUID)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get unit")
	}
	if unit.ParkingUUID == nil {
		return resp, errs.NewErrorWithDetails(errs.ErrBadRequest, "unit is disconnected of parking")
	}
	car, err := uc.ctx.Connection().CarRepository().GetByGosNumber(args.CarNumber)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get car by number")
	}

	filter := uc.ctx.Connection().SessionRepository().Filter().
		SetStatuses([]string{"active"}).
		SetCarUUIDs([]string{car.UUID.String()})
	sessions, err := uc.ctx.Connection().SessionRepository().WhereFilter(filter)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to find sessions")
	}
	if len(sessions) != 0 {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "already has sessions")
	}

	parking, err := uc.ctx.Connection().ParkingRepository().GetByUUID(unit.ParkingUUID.String())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get parking by uuid")
	}

	session := &models.Session{
		UUID:        uuid.New(),
		ParkingUUID: parking.UUID,
		CarUUID:     car.UUID,
		StartAt:     time.Now(),
		Status:      models.Active,
	}
	if err := uc.ctx.Connection().SessionRepository().Add(session); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to add session")
	}
	resp.Status = "success"
	return resp, nil
}

func (uc *SessionUseCase) Finish(args props.FinishSessionRequest, oc proto.OrderClient) (resp props.FinishSessionResp, err error) {
	// TODO: order gRPC payment connect
	// TODO: транзакции для ключевых сущностей
	car, err := uc.ctx.Connection().CarRepository().GetByGosNumber(args.CarNumber)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get car by number")
	}

	unit, err := uc.ctx.Connection().UnitRepository().GetByUUID(args.UnitUUID.String())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get unit by uuid")
	}

	// active
	filter := uc.ctx.Connection().SessionRepository().Filter().SetCarUUIDs([]string{car.UUID.String()}).SetStatuses([]string{"active"})
	sessions, err := uc.ctx.Connection().SessionRepository().WhereFilter(filter)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get session by uuid")
	}
	if len(sessions) == 0 {
		return resp, errs.NewErrorWithDetails(errs.ErrBadRequest, "you don`t have an active session")
	}

	session := sessions[0]
	parking, err := uc.ctx.Connection().ParkingRepository().GetByUUID(unit.ParkingUUID.String())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get parking by uuid")
	}

	tf, err := uc.ctx.Connection().TariffRepository().GetByUUID(parking.TariffUUID.String())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get tariff by uuid")
	}

	dur := time.Since(session.StartAt)
	if tf == nil {
		tf = &models.Tariff{}
	}

	cost, err := tariff.CalculateSessionCost(dur, *tf)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to calculate the price")
	}

	switch args.PaymentMethod {
	case props.BankCard:
		cardResp, err := oc.GetPreferredByUserUUID(context.Background(), &proto.GetPreferredCardReq{
			UserUuid: args.UserUUID,
		})
		if err != nil {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get card by uuid")
		}
		if cardResp.Card == nil {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to start session without connected bank card")
		}

		paymentResp, err := oc.CreatePayment(context.Background(), &proto.CreatePaymentReq{
			SessionUuid:   session.UUID.String(),
			PaymentMethod: string(args.PaymentMethod),
			Amount:        float32(cost),
			Description:   fmt.Sprintf("Это платеж за парковочную сессию №%s", session.UUID.String()),
			UserUuid:      args.UserUUID,
			CardUuid:      cardResp.Card.Uuid,
		})
		if err != nil {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to make payment")
		}

		if paymentResp.Payment.Status != "succeeded" {
			return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "payment failed")
		}
	}

	now := time.Now()
	updates := uc.ctx.Connection().SessionRepository().Updates().
		SetStatus(string(models.Finished)).
		SetFinishAt(&now).
		SetCost(cost)

	if err := uc.ctx.Connection().SessionRepository().Update(session.UUID, updates); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to update session")
	}
	resp.Status = "success"

	return resp, nil
}

func (uc *SessionUseCase) UpdatePaidTiming(args props.UpdateSessionPaid) error {
	updates := uc.ctx.Connection().SessionRepository().Updates()
	session, err := uc.ctx.Connection().SessionRepository().GetByUUID(args.SessionUUID.String())
	if err != nil {
		return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to find parking session")
	}

	if session == nil {
		return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to find parking session")
	}
	now := time.Now()
	updates.SetFinishAt(&now)
	if err := uc.ctx.Connection().SessionRepository().Update(args.SessionUUID, updates); err != nil {
		return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to update parking session")
	}

	return nil
}
