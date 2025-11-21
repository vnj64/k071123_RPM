package props

import "k071123/internal/services/parking_service/domain/models"

type CreateParkingReq struct {
	Name        string                  `json:"name"`
	Address     string                  `json:"address"`
	Latitude    string                  `json:"latitude"`
	Longitude   string                  `json:"longitude"`
	WorkingTime []CreateParkingSchedule `json:"working_time"`
	TotalPlaces uint                    `json:"total_places"`
	Tariff      CreateTariffReq         `json:"tariff"`
}

// TODO: days_of_week needs func int -> pq.Int64Array
type CreateParkingSchedule struct {
	DaysOfWeek []int  `gorm:"type:integer[]" json:"days_of_week"`
	OpenTime   string `json:"open_time"`
	CloseTime  string `json:"close_time"`
}

type CreateParkingResp struct {
	Parking *models.Parking `json:"parking"`
}
