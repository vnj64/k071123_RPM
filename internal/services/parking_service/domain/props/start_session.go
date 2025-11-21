package props

type StartSessionReq struct {
	CarNumber string `json:"car_number"`
	UnitUUID  string `json:"unit_uuid"`
}

type StartSessionResp struct {
	Status string `json:"status"`
}
