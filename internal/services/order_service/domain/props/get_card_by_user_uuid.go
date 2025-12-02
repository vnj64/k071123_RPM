package props

import "k071123/internal/services/order_service/domain/models"

type GetPreferredByUserUUIDReq struct {
	UserUUID string `json:"user_uuid"`
}

type GetPreferredByUserUUIDResp struct {
	Card *models.Card `json:"card"`
}
