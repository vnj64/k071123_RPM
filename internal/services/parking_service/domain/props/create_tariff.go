package props

import "k071123/internal/services/parking_service/domain/models"

type CreateTariffReq struct {
	Type models.TariffType `json:"type"`

	HasFree  *bool `json:"has_free"`
	FreeTime *int  `json:"free_time"`

	HourlyPrice     float64 `json:"hourly_price"`
	LongHourlyPrice float64 `json:"long_hourly_price"`
	DailyPrice      float64 `json:"daily_price"`

	LongHourlyStart int `json:"long_hourly_start"`
	LongHourlyEnd   int `json:"long_hourly_end"`
}

type CreateTariffResp struct {
	Tariff *models.Tariff `json:"tariff"`
}
