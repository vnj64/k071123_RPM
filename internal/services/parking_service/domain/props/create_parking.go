package props

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"k071123/internal/services/parking_service/domain/models"
	"strconv"
	"strings"
	"time"
)

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
	DaysOfWeek pq.Int64Array `gorm:"type:integer[]" json:"days_of_week"`
	OpenTime   string        `json:"open_time"`
	CloseTime  string        `json:"close_time"`
}

type CreateParkingResp struct {
	Parking *models.Parking `json:"parking"`
}

func (req *CreateParkingReq) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.Address == "" {
		return errors.New("address is required")
	}

	if err := validateCoordinate(req.Latitude); err != nil {
		return fmt.Errorf("invalid latitude: %w", err)
	}
	if err := validateCoordinate(req.Longitude); err != nil {
		return fmt.Errorf("invalid longitude: %w", err)
	}

	if req.TotalPlaces == 0 {
		return errors.New("total_places must be > 0")
	}

	if err := req.Tariff.Validate(); err != nil {
		return fmt.Errorf("invalid tariff: %w", err)
	}

	if len(req.WorkingTime) == 0 {
		return errors.New("working_time is required")
	}

	for i, wt := range req.WorkingTime {
		if err := wt.Validate(); err != nil {
			return fmt.Errorf("working_time[%d]: %w", i, err)
		}
	}

	return nil
}

func (t *CreateTariffReq) Validate() error {
	if t.Type == "" {
		return errors.New("tariff.type is required")
	}
	if t.HourlyPrice < 0 {
		return errors.New("hourly_price must be >= 0")
	}
	if t.DailyPrice < 0 {
		return errors.New("daily_price must be >= 0")
	}
	return nil
}

func (s *CreateParkingSchedule) Validate() error {
	if len(s.DaysOfWeek) == 0 {
		return errors.New("days_of_week must not be empty")
	}

	for _, day := range s.DaysOfWeek {
		if day < 1 || day > 7 {
			return fmt.Errorf("invalid day_of_week %d: allowed 1–7", day)
		}
	}

	if s.OpenTime == "" || s.CloseTime == "" {
		return errors.New("open_time and close_time are required")
	}

	open, err := parseHHMM(s.OpenTime)
	if err != nil {
		return fmt.Errorf("invalid open_time: %w", err)
	}

	closing, err := parseHHMM(s.CloseTime)
	if err != nil {
		return fmt.Errorf("invalid close_time: %w", err)
	}

	if !closing.After(open) {
		return errors.New("close_time must be after open_time")
	}

	return nil
}

func validateCoordinate(value string) error {
	if value == "" {
		return errors.New("value is empty")
	}
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New("must be a valid float")
	}
	return nil
}

func parseHHMM(s string) (time.Time, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return time.Time{}, errors.New("format must be HH:MM")
	}

	hour, err1 := strconv.Atoi(parts[0])
	min, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return time.Time{}, errors.New("invalid numeric values")
	}

	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		return time.Time{}, errors.New("time out of range")
	}

	// Возвращаем "время внутри дня"
	return time.Date(2000, 1, 1, hour, min, 0, 0, time.UTC), nil
}
