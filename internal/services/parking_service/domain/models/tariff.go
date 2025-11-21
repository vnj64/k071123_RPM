package models

import (
	"github.com/google/uuid"
	"k071123/pkg/timestamps"
)

type Tariff struct {
	Timestamps timestamps.Timestamps `gorm:"embedded;" json:"-" swaggerignore:"true"`
	UUID       uuid.UUID             `json:"uuid" gorm:"primaryKey"`

	// TODO: подумать над необходиомостью
	Type TariffType `json:"type"`

	HasFree  *bool `json:"has_free"`
	FreeTime *int  `json:"free_time"`

	HourlyPrice     float64 `json:"hourly_price"`
	LongHourlyPrice float64 `json:"long_hourly_price"`
	DailyPrice      float64 `json:"daily_price"`

	LongHourlyStart int `json:"long_hourly_start"`
	LongHourlyEnd   int `json:"long_hourly_end"`
}

type TariffType string

const (
	Free       TariffType = "free"
	Hourly     TariffType = "hourly"
	LongHourly TariffType = "long_hourly"
	Daily      TariffType = "daily"
)
