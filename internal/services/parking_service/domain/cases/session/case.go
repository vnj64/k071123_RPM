package session

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	proto2 "k071123/internal/services/notification_service/contracts/pkg/proto"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/cases/tariff"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	proto3 "k071123/internal/services/user_service/contracts/pkg/proto"
	"k071123/internal/utils/errs"
	"time"
)

type SessionUseCase struct {
	ctx domain.Context
	oc  proto.OrderClient
	nc  proto2.NotificationClient
	uc  proto3.UserClient
}

func NewSessionUseCase(
	ctx domain.Context,
	oc proto.OrderClient,
	nc proto2.NotificationClient,
	uc proto3.UserClient,
) *SessionUseCase {
	return &SessionUseCase{ctx: ctx, oc: oc, nc: nc, uc: uc}
}

// TODO: сделать Search по Sessions

func (uc *SessionUseCase) Start(args props.StartSessionReq, oc proto.OrderClient) (resp props.StartSessionResp, err error) {
	// unit parking uuid validation
	log := uc.ctx.Services().Logger().WithField("SessionUseCase", "Start")
	tx, err := uc.ctx.Connection().Begin()
	if err != nil {
		log.Errorf("begin transaction error: %s", err.Error())
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to start transaction")
	}
	defer tx.Rollback()

	cardResp, err := oc.GetPreferredByUserUUID(context.Background(), &proto.GetPreferredCardReq{
		UserUuid: args.UserUUID,
	})
	if err != nil {
		log.Errorf("grpc error, cannot get card: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get user card from order service")
	}
	if cardResp.Card == nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "cannot start session without connected card")
	}

	unit, err := tx.UnitRepository().GetByUUID(args.UnitUUID)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get unit")
	}
	if unit.ParkingUUID == nil {
		return resp, errs.NewErrorWithDetails(errs.ErrBadRequest, "unit is disconnected of parking")
	}
	car, err := tx.CarRepository().GetByGosNumber(args.CarNumber)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to get car by number")
	}

	filter := tx.SessionRepository().Filter().
		SetStatuses([]string{"active"}).
		SetCarUUIDs([]string{car.UUID.String()})
	sessions, err := tx.SessionRepository().WhereFilter(filter)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to find sessions")
	}
	if len(sessions) != 0 {
		log.Infof("found sessions: %v", sessions)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "already has sessions")
	}

	parking, err := tx.ParkingRepository().GetByUUID(unit.ParkingUUID.String())
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
	if err := tx.SessionRepository().Add(session); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to add session")
	}
	if err := tx.Commit(); err != nil {
		log.Errorf("commit transaction error: %s", err.Error())
		return resp, err
	}
	resp.Status = "success"
	return resp, nil
}

func (uc *SessionUseCase) Finish(args props.FinishSessionRequest, oc proto.OrderClient) (resp props.FinishSessionResp, err error) {
	log := uc.ctx.Services().Logger().WithField("SessionUseCase", "Finish")
	tx, err := uc.ctx.Connection().Begin()
	if err != nil {
		log.Errorf("begin transaction error: %s", err.Error())
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to start transaction")
	}
	defer tx.Rollback()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Infof("user uuid: %s", args.UserUUID)
	user, err := uc.uc.GetUserByUUID(ctx, &proto3.GetUserReq{
		Uuid: args.UserUUID,
	})
	if user == nil {
		log.Errorf("cannot find user")
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to find user")
	}

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

	if err := tx.SessionRepository().Update(session.UUID, updates); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to update session")
	}
	if err := tx.Commit(); err != nil {
		return resp, err
	}

	// user service grpc
	// OUTBOX
	_, err = uc.nc.SendEmail(ctx, &proto2.SendEmailReq{
		Data:    "Your parking session successfully finished",
		Subject: "Parking Session",
		To:      []string{user.Email},
	})
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to send email")
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
