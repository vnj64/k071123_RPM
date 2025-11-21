package tariff

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
	"time"
)

type TariffUseCase struct {
	ctx domain.Context
}

func NewTariffUseCase(ctx domain.Context) *TariffUseCase {
	return &TariffUseCase{ctx: ctx}
}

func (uc *TariffUseCase) Create(args props.CreateTariffReq) (resp props.CreateTariffResp, err error) {
	tariff := &models.Tariff{
		UUID:            uuid.New(),
		HourlyPrice:     args.HourlyPrice,
		LongHourlyPrice: args.LongHourlyPrice,
		DailyPrice:      args.DailyPrice,
		LongHourlyStart: args.LongHourlyStart,
		LongHourlyEnd:   args.LongHourlyEnd,
		Type:            args.Type,
	}
	if args.HasFree != nil {
		tariff.HasFree = args.HasFree
	}
	if args.FreeTime != nil {
		tariff.FreeTime = args.FreeTime
	}

	if err := uc.ctx.Connection().TariffRepository().Add(&models.Tariff{
		UUID: uuid.New(),
	}); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "database error")
	}
	resp.Tariff = tariff
	return resp, nil
}

func CalculateSessionCost(duration time.Duration, tariff models.Tariff) (float64, error) {
	cost := 0.0
	durationHours := duration.Hours()
	if tariff.HasFree != nil {
		// бесплатно
		if duration.Minutes() < float64(*tariff.FreeTime) {
			return 0, nil
		}
	}
	// меньше часа
	if duration.Minutes() < float64(time.Duration(60)*time.Minute) {
		cost = tariff.HourlyPrice
		return cost, nil
	}
	// при условии, что 60 < duration < lhs
	lhs := time.Duration(tariff.LongHourlyStart) * time.Minute // LongHourlyStart в минутах
	if duration.Minutes() < lhs.Minutes() && duration.Minutes() > 60 {
		cost = tariff.HourlyPrice * durationHours
		return cost, nil
	}
	// lhs < duration < lhe
	lhe := time.Duration(tariff.LongHourlyEnd) * time.Minute
	if duration.Minutes() < lhe.Minutes() && duration.Minutes() > 60 {
		cost = tariff.LongHourlyPrice
		return cost, nil
	}
	// lhe < duration
	if lhe.Minutes() < duration.Minutes() {
		if duration.Hours() >= 48 {
			cost = tariff.DailyPrice * (duration.Hours() / 24)
			return cost, nil
		}
		cost = tariff.DailyPrice
		return cost, nil
	}
	return 0, nil
}
