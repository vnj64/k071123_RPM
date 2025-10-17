package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"k071123/internal/services/notification_service/contracts/pkg/proto"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/models/user_status"
	"k071123/internal/services/user_service/domain/props"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	auth "k071123/internal/utils/jwt_helpers"
	"k071123/pkg/timestamps"
	"math/rand"
	"strconv"
	"time"
)

type AuthUseCase struct {
	ctx          domain.Context
	notifyClient proto.NotificationClient
}

func NewAuthUseCase(ctx domain.Context, nC proto.NotificationClient) *AuthUseCase {
	return &AuthUseCase{ctx: ctx, notifyClient: nC}
}

func (uc *AuthUseCase) SendCode(args props.SendCodeReq) (resp props.SendCodeResponse, err error) {
	code := generateRandomToken()
	message := fmt.Sprintf("Your confirmation code for PRK Project: %s", code)
	subject := "Your confirmation code"
	response, err := uc.notifyClient.SendEmail(context.Background(), &proto.SendEmailReq{
		Data:    message,
		Subject: subject,
		To: []string{
			args.Email,
		},
	})
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, err.Error())
	}
	if response.Response == "failed" {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to send email")
	}

	if err := uc.ctx.Connection().VerificationCode().Add(&models.VerificationCode{
		UUID:      uuid.New(),
		Email:     args.Email,
		Code:      code,
		Used:      false,
		CreatedAt: time.Now(),
	}); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, err.Error())
	}

	now := time.Now()
	user := &models.User{
		UUID:   uuid.New(),
		Status: user_status.Inactive,
		Email:  args.Email,
		Role:   permissions.Default,
		Timestamps: timestamps.Timestamps{
			CreatedAt: now,
			UpdatedAt: &now,
		},
	}
	if err := uc.ctx.Connection().User().Add(user); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}

	return props.SendCodeResponse{
		Status: "success",
	}, nil
}

func (uc *AuthUseCase) ConfirmCode(args props.ConfirmCodeReq) (resp props.ConfirmCodeResp, err error) {
	verificationCode, err := uc.ctx.Connection().VerificationCode().GetByCode(args.Code)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}
	if verificationCode.Email != args.Email {
		return resp, errs.NewErrorWithDetails(errs.ErrForbidden, "you are now owner of this code")
	}
	if args.Code != verificationCode.Code {
		return resp, errs.NewErrorWithDetails(errs.ErrForbidden, "wrong code")
	}

	user, err := uc.ctx.Connection().User().GetByEmail(args.Email)
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}

	acToken, err := auth.GenerateAuthToken(user.UUID.String(), user.Role, uc.ctx.Services().Config())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "cannot generate auth token")
	}

	refToken, err := auth.GenerateRefreshToken(user.UUID, uc.ctx.Services().Config())
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "cannot generate refresh token")
	}
	resp.AccessToken = acToken
	resp.RefreshToken = refToken.RefreshTokenUUID.String()

	return resp, nil
}

func generateRandomToken() string {
	var code string
	for i := 0; i < 5; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	return code
}
