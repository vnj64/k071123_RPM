package session

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/cases/tariff"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
	"log"
	"time"
)

type SessionUseCase struct {
	ctx domain.Context
}

func NewSessionUseCase(ctx domain.Context) *SessionUseCase {
	return &SessionUseCase{ctx: ctx}
}

func (uc *SessionUseCase) Start(args props.StartSessionReq) (resp props.StartSessionResp, err error) {
	// unit parking uuid validation
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

func (uc *SessionUseCase) Finish(args props.FinishSessionRequest) (resp props.FinishSessionResp, err error) {
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
	// TODO: order gRPC payment connect
	// AutoPay(token string, amount float64) error
	log.Printf("cost %v", cost)
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
	if err := uc.ctx.Connection().SessionRepository().Update(uc.ctx.Connection().DB(), args.SessionUUID, updates); err != nil {
		return errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to update parking session")
	}

	return nil
}
