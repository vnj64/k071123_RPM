package usecases

import (
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	mocks2 "k071123/internal/services/parking_service/domain/mocks"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/services/parking_service/tests/mocks"
	"k071123/tools/logger"
	"testing"

	"k071123/internal/services/parking_service/domain/cases/car"
)

func TestCarUseCase_CreateCar(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	contextMock := mocks2.NewContextMock(mc)
	servicesMock := mocks2.NewServicesMock(mc)
	connMock := mocks2.NewConnectionMock(mc)

	carRepoMock := mocks.NewCarRepositoryMock(mc)

	contextMock.ServicesMock.Return(servicesMock)
	contextMock.ConnectionMock.Return(connMock)

	connMock.CarRepositoryMock.Return(carRepoMock)

	realLogrus := logrus.New()
	fakeLog := &logger.Logger{
		Logger: realLogrus,
	}
	servicesMock.LoggerMock.Return(fakeLog)

	tests := []struct {
		name       string
		args       props.CreateCarReq
		wantErr    assert.ErrorAssertionFunc
		mockExpect func()
	}{
		{
			name: "Success Create",
			args: props.CreateCarReq{
				GosNumber: "B412HB716",
				UserUUID:  uuid.New().String(),
			},
			wantErr: assert.NoError,
			mockExpect: func() {
				carRepoMock.AddMock.
					Return(nil)
			},
		},
		{
			name: "Validation Error (empty gosnumber)",
			args: props.CreateCarReq{
				GosNumber: "",
				UserUUID:  uuid.New().String(),
			},
			wantErr:    assert.Error,
			mockExpect: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			uc := car.NewCarUseCase(contextMock)

			_, err := uc.CreateCar(tt.args)
			tt.wantErr(t, err)
		})
	}
}
