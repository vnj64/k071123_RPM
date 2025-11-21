package car

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
)

type CarUseCase struct {
	ctx domain.Context
}

func NewCarUseCase(ctx domain.Context) *CarUseCase {
	return &CarUseCase{ctx: ctx}
}

func (uc *CarUseCase) CreateCar(args props.CreateCarReq) (resp props.CreateCarResp, err error) {
	if err := args.Validate(); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, err.Error())
	}

	car := &models.Car{
		UUID:      uuid.New(),
		GosNumber: args.GosNumber,
		IsActive:  true,
	}
	if err := uc.ctx.Connection().CarRepository().Add(car); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}
	resp.Car = car

	return resp, nil
}
