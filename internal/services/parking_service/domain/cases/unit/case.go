package unit

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
)

type UnitUseCase struct {
	ctx domain.Context
}

// TODO: CRUD UNIT

func NewUnitUseCase(ctx domain.Context) *UnitUseCase {
	return &UnitUseCase{ctx: ctx}
}

func (uc *UnitUseCase) CreateUnit(args props.CreateUniqReq) (resp props.CreateUnitResp, err error) {
	log := uc.ctx.Services().Logger().WithField("UnitUseCase", "CreateUnit")
	if err := args.Validate(); err != nil {
		log.Errorf("validation error on unit: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, err.Error())
	}
	unit := &models.Unit{
		UUID:          uuid.New(),
		Status:        args.Status,
		NetworkStatus: args.NetworkStatus,
		Direction:     args.Direction,
	}
	if args.Code != nil {
		unit.Code = args.Code
	}
	if args.QrLink != nil {
		unit.QrLink = args.QrLink
	}
	if args.ParkingUUID != nil {
		unit.ParkingUUID = args.ParkingUUID
	}

	if err := uc.ctx.Connection().UnitRepository().Add(unit); err != nil {
		log.Errorf("unable to add unit: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to insert unit")
	}

	resp.Unit = unit

	return resp, nil
}

func (uc *UnitUseCase) GetByUUID(args props.GetUnitByUUID) (resp props.GetUnitByUUIDResp, err error) {
	log := uc.ctx.Services().Logger().WithField("UnitUseCase", "GetByUUID").WithFields(logrus.Fields{})
	if args.UUID == uuid.Nil {
		log.Errorf("uuid is nil on request")
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "uuid is nil on request")
	}

	unit, err := uc.ctx.Connection().UnitRepository().GetByUUID(args.UUID.String())
	if err != nil {
		log.Errorf("unable to get unit by uuid: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unit not found")
	}
	resp.Unit = unit
	return resp, nil
}
