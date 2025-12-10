package props

import "k071123/internal/services/user_service/domain/models"

type GetUserByUUIDReq struct {
	UUID string `json:"uuid"`
}

type GetUserByUUIDResp struct {
	User *models.User `json:"user"`
}
