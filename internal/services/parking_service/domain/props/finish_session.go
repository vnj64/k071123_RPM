package props

import "github.com/google/uuid"

type FinishSessionRequest struct {
	CarNumber string    `json:"session_uuid"`
	UnitUUID  uuid.UUID `json:"unit_uuid"`
}

type FinishSessionResp struct {
	Status string `json:"status"`
}
