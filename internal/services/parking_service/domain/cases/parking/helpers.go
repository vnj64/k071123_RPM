package parking

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
)

func ParkingScheduleParser(args props.CreateParkingSchedule, parkingUUID uuid.UUID) models.ParkingSchedule {
	var days64type []int64
	for i := 0; i < len(args.DaysOfWeek); i++ {
		days64type = append(days64type, int64(i))
	}

	return models.ParkingSchedule{
		UUID:        uuid.New(),
		DaysOfWeek:  args.DaysOfWeek,
		OpenTime:    args.OpenTime,
		CloseTime:   args.CloseTime,
		ParkingUUID: parkingUUID,
	}
}
