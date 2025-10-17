package user

import (
	"github.com/google/uuid"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/props"
	"k071123/internal/utils/errs"
	"k071123/pkg/timestamps"
	"time"
)

type UserUseCase struct {
	ctx domain.Context
}

func NewUserUseCase(ctx domain.Context) *UserUseCase {
	return &UserUseCase{ctx: ctx}
}

func (uc *UserUseCase) CreateUser(args *props.CreateUserReq) (resp props.CreateUserResp, err error) {
	if err := args.Validate(); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrUnprocessableEntity, err.Error())
	}

	now := time.Now()
	model := &models.User{
		UUID:        uuid.New(),
		Status:      args.Status,
		PhoneNumber: &args.PhoneNumber,
		Timestamps: timestamps.Timestamps{
			CreatedAt: now,
			UpdatedAt: &now,
		},
	}

	if args.FirstName != nil {
		model.FirstName = args.FirstName
	}
	if args.SecondName != nil {
		model.SecondName = args.SecondName
	}
	if args.BirthDate != nil {
		model.BirthDate = args.BirthDate
	}

	if err := uc.ctx.Connection().User().Add(model); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, err.Error())
	}
	resp.User = model

	return resp, nil
}
