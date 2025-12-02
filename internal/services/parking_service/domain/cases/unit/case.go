package unit

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
)

type UnitUseCase struct {
	ctx domain.Context
}

func NewUnitUseCase(ctx domain.Context) *UnitUseCase {
	return &UnitUseCase{ctx: ctx}
}

func (uc *UnitUseCase) Create(args props.CreateUniqReq) (resp props.CreateUnitResp, err error) {
	uc.ctx.Services().Logger().Info("create unit start")
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
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to insert unit")
	}

	resp.Unit = unit
	return resp, nil
}
