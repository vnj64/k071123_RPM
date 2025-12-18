package usecases

import (
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	proto3 "k071123/internal/services/notification_service/contracts/pkg/proto"
	mocks3 "k071123/internal/services/notification_service/contracts/pkg/proto/mocks"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	orderMocks "k071123/internal/services/order_service/contracts/pkg/proto/mocks"
	"k071123/internal/services/parking_service/domain/cases/session"
	mocks2 "k071123/internal/services/parking_service/domain/mocks"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	mocks5 "k071123/internal/services/parking_service/domain/repositories/mocks"
	"k071123/internal/services/parking_service/tests/mocks"
	proto2 "k071123/internal/services/user_service/contracts/pkg/proto"
	mocks4 "k071123/internal/services/user_service/contracts/pkg/proto/mocks"
	"k071123/tools/logger"
	"testing"
	"time"
)

func TestSessionUseCase_StartSession(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := mocks2.NewContextMock(mc)
	services := mocks2.NewServicesMock(mc)
	conn := mocks2.NewConnectionMock(mc)
	tx := mocks2.NewTransactionalConnectionMock(mc)

	orderClient := orderMocks.NewOrderClientMock(mc)
	notificationClient := mocks3.NewNotificationClientMock(mc)
	userClient := mocks4.NewUserClientMock(mc)

	sessionRepo := mocks.NewSessionRepositoryMock(mc)
	unitRepo := mocks.NewUnitRepositoryMock(mc)
	parkingRepo := mocks.NewParkingRepositoryMock(mc)
	carRepo := mocks.NewCarRepositoryMock(mc)
	sessionFilter := mocks5.NewSessionFilterMock(mc)

	log := &logger.Logger{Logger: logrus.New()}
	services.LoggerMock.Return(log)

	ctx.ServicesMock.Return(services)
	ctx.ConnectionMock.Return(conn)

	userUUID := uuid.New()
	car := &models.Car{
		UUID:      uuid.New(),
		GosNumber: "B412HB716",
		UserUUID:  userUUID,
		IsActive:  true,
	}

	unitUUID := uuid.New()
	parkingUUID := uuid.New()

	tests := []struct {
		name       string
		args       props.StartSessionReq
		wantErr    assert.ErrorAssertionFunc
		mockExpect func()
	}{
		{
			name: "Success Start",
			args: props.StartSessionReq{
				CarNumber: car.GosNumber,
				UnitUUID:  unitUUID.String(),
				UserUUID:  userUUID.String(),
			},
			wantErr: assert.NoError,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				orderClient.GetPreferredByUserUUIDMock.Return(
					&proto.GetPreferredCardResp{
						Card: &proto.Card{
							Uuid: uuid.New().String(),
						},
					},
					nil,
				)

				tx.UnitRepositoryMock.Return(unitRepo)
				tx.CarRepositoryMock.Return(carRepo)
				tx.SessionRepositoryMock.Return(sessionRepo)
				tx.ParkingRepositoryMock.Return(parkingRepo)

				unitRepo.GetByUUIDMock.Return(&models.Unit{
					UUID:        unitUUID,
					ParkingUUID: &parkingUUID,
				}, nil)

				carRepo.GetByGosNumberMock.Return(car, nil)

				sessionRepo.FilterMock.Return(sessionFilter)
				sessionFilter.SetStatusesMock.Return(sessionFilter)
				sessionFilter.SetCarUUIDsMock.Return(sessionFilter)
				sessionRepo.WhereFilterMock.Return([]models.Session{}, nil)

				parkingRepo.GetByUUIDMock.Return(&models.Parking{
					UUID: parkingUUID,
				}, nil)

				sessionRepo.AddMock.Return(nil)

				tx.CommitMock.Return(nil)
			},
		},
		{
			name: "Validation Error",
			args: props.StartSessionReq{
				CarNumber: "",
				UnitUUID:  unitUUID.String(),
				UserUUID:  userUUID.String(),
			},
			wantErr: assert.Error,
			mockExpect: func() {
			},
		},
		{
			name: "Unavailable Bank Card",
			args: props.StartSessionReq{
				CarNumber: car.GosNumber,
				UnitUUID:  unitUUID.String(),
				UserUUID:  userUUID.String(),
			},
			wantErr: assert.Error,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				orderClient.GetPreferredByUserUUIDMock.Return(
					&proto.GetPreferredCardResp{
						Card: nil,
					},
					nil,
				)
			},
		},
		{
			name: "Already Active Session",
			args: props.StartSessionReq{
				CarNumber: car.GosNumber,
				UnitUUID:  unitUUID.String(),
				UserUUID:  userUUID.String(),
			},
			wantErr: assert.Error,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				orderClient.GetPreferredByUserUUIDMock.Return(
					&proto.GetPreferredCardResp{
						Card: &proto.Card{
							Uuid: uuid.New().String(),
						},
					},
					nil,
				)
				tx.UnitRepositoryMock.Return(unitRepo)
				tx.CarRepositoryMock.Return(carRepo)
				tx.SessionRepositoryMock.Return(sessionRepo)
				tx.ParkingRepositoryMock.Return(parkingRepo)

				unitRepo.GetByUUIDMock.Return(&models.Unit{
					UUID:        unitUUID,
					ParkingUUID: &parkingUUID,
				}, nil)

				carRepo.GetByGosNumberMock.Return(car, nil)

				sessionRepo.FilterMock.Return(sessionFilter)
				sessionFilter.SetStatusesMock.Return(sessionFilter)
				sessionFilter.SetCarUUIDsMock.Return(sessionFilter)
				sessionRepo.WhereFilterMock.Return([]models.Session{
					{
						UUID: uuid.New(),
					},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			uc := session.NewSessionUseCase(
				ctx,
				orderClient,
				notificationClient,
				userClient,
			)

			_, err := uc.Start(tt.args)
			tt.wantErr(t, err)
		})
	}
}

func TestSessionUseCase_FinishSession(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := mocks2.NewContextMock(mc)
	services := mocks2.NewServicesMock(mc)
	conn := mocks2.NewConnectionMock(mc)
	tx := mocks2.NewTransactionalConnectionMock(mc)

	orderClient := orderMocks.NewOrderClientMock(mc)
	notificationClient := mocks3.NewNotificationClientMock(mc)
	userClient := mocks4.NewUserClientMock(mc)

	sessionRepo := mocks.NewSessionRepositoryMock(mc)
	unitRepo := mocks.NewUnitRepositoryMock(mc)
	parkingRepo := mocks.NewParkingRepositoryMock(mc)
	carRepo := mocks.NewCarRepositoryMock(mc)
	tariffRepo := mocks.NewTariffRepositoryMock(mc)

	sessionFilter := mocks5.NewSessionFilterMock(mc)
	sessionUpdates := mocks5.NewSessionUpdatesMock(mc)

	log := &logger.Logger{Logger: logrus.New()}
	services.LoggerMock.Return(log)

	ctx.ServicesMock.Return(services)
	ctx.ConnectionMock.Return(conn)

	var (
		parkingUUID = uuid.New()
		userUUID    = uuid.New()

		unit = &models.Unit{
			UUID:        uuid.New(),
			Code:        ptrStr("123456"),
			ParkingUUID: &parkingUUID,
		}
		car = &models.Car{
			UUID:      uuid.New(),
			GosNumber: "B412HB716",
			IsActive:  true,
			UserUUID:  userUUID,
		}
		tariff = &models.Tariff{
			UUID:            uuid.New(),
			HasFree:         ptrBool(true),
			FreeTime:        ptrInt(10),
			HourlyPrice:     100,
			LongHourlyPrice: 10,
			LongHourlyStart: 300,
			LongHourlyEnd:   600,
			DailyPrice:      1000,
		}
		parking = &models.Parking{
			UUID:       parkingUUID,
			TariffUUID: tariff.UUID,
		}
	)

	tests := []struct {
		name       string
		args       props.FinishSessionRequest
		wantErr    assert.ErrorAssertionFunc
		mockExpect func()
	}{
		{
			name: "Success Finish",
			args: props.FinishSessionRequest{
				CarNumber:     car.GosNumber,
				UnitUUID:      unit.UUID,
				UserUUID:      userUUID.String(),
				PaymentMethod: "bank_card",
			},
			wantErr: assert.NoError,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				userClient.GetUserByUUIDMock.Return(&proto2.GetUserResp{
					Uuid:      userUUID.String(),
					FirstName: "Kamil",
					Email:     "vnj644@gmail.com",
				}, nil)

				tx.CarRepositoryMock.Return(carRepo)
				tx.UnitRepositoryMock.Return(unitRepo)
				tx.SessionRepositoryMock.Return(sessionRepo)
				tx.ParkingRepositoryMock.Return(parkingRepo)
				tx.TariffRepositoryMock.Return(tariffRepo)

				carRepo.GetByGosNumberMock.Return(car, nil)

				unitRepo.GetByUUIDMock.Return(unit, nil)

				sessionRepo.FilterMock.Return(sessionFilter)
				sessionFilter.SetCarUUIDsMock.Return(sessionFilter)
				sessionFilter.SetStatusesMock.Return(sessionFilter)
				sessionRepo.WhereFilterMock.Return([]models.Session{
					{
						UUID:    uuid.New(),
						CarUUID: car.UUID,
						StartAt: time.Now().Add(-1 * time.Hour),
						Status:  "active",
					},
				}, nil)

				parkingRepo.GetByUUIDMock.Return(parking, nil)

				tariffRepo.GetByUUIDMock.Return(tariff, nil)

				orderClient.GetPreferredByUserUUIDMock.Return(&proto.GetPreferredCardResp{
					Card: &proto.Card{
						Uuid: uuid.New().String(),
					},
				}, nil)

				orderClient.CreatePaymentMock.Return(&proto.CreatePaymentResp{
					Payment: &proto.Payment{
						Uuid:   uuid.New().String(),
						Status: "succeeded",
					},
				}, nil)

				sessionRepo.UpdatesMock.Return(sessionUpdates)
				sessionUpdates.SetStatusMock.Return(sessionUpdates)
				sessionUpdates.SetFinishAtMock.Return(sessionUpdates)
				sessionUpdates.SetCostMock.Return(sessionUpdates)

				sessionRepo.UpdateMock.Return(nil)
				tx.CommitMock.Return(nil)

				notificationClient.SendEmailMock.Return(&proto3.SendEmailResp{Response: "success"}, nil)
			},
		},
		{
			name: "Validation Error",
			args: props.FinishSessionRequest{
				CarNumber:     "",
				UnitUUID:      unit.UUID,
				UserUUID:      userUUID.String(),
				PaymentMethod: "bank_card",
			},
			wantErr: assert.Error,
			mockExpect: func() {

			},
		},
		{
			name: "User is Empty",
			args: props.FinishSessionRequest{
				CarNumber:     car.GosNumber,
				UnitUUID:      unit.UUID,
				UserUUID:      userUUID.String(),
				PaymentMethod: "bank_card",
			},
			wantErr: assert.Error,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				userClient.GetUserByUUIDMock.Return(nil, nil)
			},
		},
		{
			name: "No active sessions",
			args: props.FinishSessionRequest{
				CarNumber:     car.GosNumber,
				UnitUUID:      unit.UUID,
				UserUUID:      userUUID.String(),
				PaymentMethod: "bank_card",
			},
			wantErr: assert.Error,
			mockExpect: func() {
				conn.BeginMock.Return(tx, nil)

				userClient.GetUserByUUIDMock.Return(&proto2.GetUserResp{
					Uuid:      userUUID.String(),
					FirstName: "Kamil",
					Email:     "vnj644@gmail.com",
				}, nil)

				tx.CarRepositoryMock.Return(carRepo)
				tx.UnitRepositoryMock.Return(unitRepo)
				tx.SessionRepositoryMock.Return(sessionRepo)
				tx.ParkingRepositoryMock.Return(parkingRepo)
				tx.TariffRepositoryMock.Return(tariffRepo)

				carRepo.GetByGosNumberMock.Return(car, nil)

				unitRepo.GetByUUIDMock.Return(unit, nil)

				sessionRepo.FilterMock.Return(sessionFilter)
				sessionFilter.SetCarUUIDsMock.Return(sessionFilter)
				sessionFilter.SetStatusesMock.Return(sessionFilter)
				sessionRepo.WhereFilterMock.Return([]models.Session{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			uc := session.NewSessionUseCase(
				ctx,
				orderClient,
				notificationClient,
				userClient,
			)

			_, err := uc.Finish(tt.args)
			tt.wantErr(t, err)
		})
	}
}
