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

// TODO: CRUD + GetByGosNumber

func NewCarUseCase(ctx domain.Context) *CarUseCase {
	return &CarUseCase{ctx: ctx}
}

func (uc *CarUseCase) CreateCar(args props.CreateCarReq) (resp props.CreateCarResp, err error) {
	log := uc.ctx.Services().Logger().WithField("CarUseCase", "CreateCar")
	if err := args.Validate(); err != nil {
		log.Errorf("validate input error: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, err.Error())
	}

	car := &models.Car{
		UUID:      uuid.New(),
		GosNumber: args.GosNumber,
		IsActive:  true,
		UserUUID:  uuid.MustParse(args.UserUUID),
	}
	if err := uc.ctx.Connection().CarRepository().Add(car); err != nil {
		log.Errorf("failed to add car to database: %v", err)
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}
	resp.Car = car

	return resp, nil
}
