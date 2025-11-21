package workers

import (
	"github.com/google/uuid"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/models/user_status"
	"k071123/internal/shared/permissions"
	"log"
)

func CreateAdminWorker(ctx domain.Context) {
	admin := &models.User{
		UUID:   uuid.New(),
		Role:   permissions.SuperAdmin,
		Status: user_status.Active,
		Email:  ctx.Services().Config().AdminLogin(),
	}

	user, err := ctx.Connection().User().GetByEmail(admin.Email)
	if err != nil {
		log.Printf("error getting user by email: %v", err)
	}
	if user == nil {
		if err := ctx.Connection().User().Add(admin); err != nil {
			log.Printf("user creating error: %v [REPO]", err.Error())
		}
	}
}
