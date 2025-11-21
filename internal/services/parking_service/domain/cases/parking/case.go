package parking

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/models/parking_statuses"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/utils/errs"
	"log"
)

type ParkingUseCase struct {
	ctx domain.Context
}

func NewParkingUseCase(ctx domain.Context) *ParkingUseCase {
	return &ParkingUseCase{ctx: ctx}
}

func (uc *ParkingUseCase) CreateParking(args props.CreateParkingReq) (resp props.CreateParkingResp, err error) {
	tx, err := uc.ctx.Connection().Begin()
	if err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "unable to start transaction")
	}
	defer tx.Rollback()

	parking := &models.Parking{
		UUID:        uuid.New(),
		Name:        args.Name,
		Latitude:    args.Latitude,
		Longitude:   args.Longitude,
		Address:     args.Address,
		TotalPlaces: args.TotalPlaces,
		Status:      parking_statuses.Active,
	}

	// TODO: функционал парсинга времени вынести в отдельную функцию
	log.Print(len(args.WorkingTime))

	var parkingScheduleList []models.ParkingSchedule
	for i := 0; i < len(args.WorkingTime); i++ {
		workingTime := ParkingScheduleParser(args.WorkingTime[i], parking.UUID)
		parkingScheduleList = append(parkingScheduleList, workingTime)
	}
	//parking.WorkingTime = parkingScheduleList

	log.Printf("%+v", len(parking.WorkingTime))
	log.Printf("%+v", len(parkingScheduleList))

	tariff := &models.Tariff{
		UUID:            uuid.New(),
		Type:            args.Tariff.Type,
		HasFree:         args.Tariff.HasFree,
		FreeTime:        args.Tariff.FreeTime,
		HourlyPrice:     args.Tariff.HourlyPrice,
		LongHourlyPrice: args.Tariff.LongHourlyPrice,
		DailyPrice:      args.Tariff.DailyPrice,
		LongHourlyEnd:   args.Tariff.LongHourlyEnd,
		LongHourlyStart: args.Tariff.LongHourlyStart,
	}

	parking.TariffUUID = tariff.UUID
	if err := uc.ctx.Connection().TariffRepository().Add(tariff); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to create tariff")
	}

	if err := uc.ctx.Connection().ParkingRepository().Add(parking); err != nil {
		return resp, errs.NewErrorWithDetails(errs.ErrInternalServerError, "failed to create parking")
	}
	resp.Parking = parking

	err = tx.Commit()
	if err != nil {
		return props.CreateParkingResp{}, err
	}

	return resp, nil
}
