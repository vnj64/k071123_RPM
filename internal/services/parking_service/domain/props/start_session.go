package props

import "errors"

type StartSessionReq struct {
	CarNumber string `json:"car_number"`
	UnitUUID  string `json:"unit_uuid"`
	UserUUID  string `json:"-"`
}

func (r *StartSessionReq) Validate() error {
	if r.CarNumber == "" {
		return errors.New("car number is required")
	}
	if r.UnitUUID == "" {
		return errors.New("unit_uuid is required")
	}
	if r.UserUUID == "" {
		return errors.New("user_uuid is required")
	}
	return nil
}

type StartSessionResp struct {
	Status string `json:"status"`
}
