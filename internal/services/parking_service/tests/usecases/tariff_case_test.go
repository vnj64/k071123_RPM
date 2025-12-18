package usecases

import (
	"github.com/gojuno/minimock/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k071123/internal/services/parking_service/domain/cases/tariff"
	mocks2 "k071123/internal/services/parking_service/domain/mocks"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/services/parking_service/tests/mocks"
	"k071123/tools/logger"
	"testing"
)

func TestTariffUseCase_CreateTariff(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	contextMock := mocks2.NewContextMock(mc)
	servicesMock := mocks2.NewServicesMock(mc)
	connMock := mocks2.NewConnectionMock(mc)

	tariffRepoMock := mocks.NewTariffRepositoryMock(mc)

	contextMock.ServicesMock.Return(servicesMock)
	contextMock.ConnectionMock.Return(connMock)

	connMock.TariffRepositoryMock.Return(tariffRepoMock)

	realLogrus := logrus.New()
	fakeLog := &logger.Logger{
		Logger: realLogrus,
	}
	servicesMock.LoggerMock.Return(fakeLog)

	tests := []struct {
		name       string
		args       props.CreateTariffReq
		wantErr    assert.ErrorAssertionFunc
		mockExpect func()
	}{
		{
			name: "Success Create",
			args: props.CreateTariffReq{
				Type:            models.Free,
				HasFree:         ptrBool(false),
				HourlyPrice:     100,
				DailyPrice:      1000,
				LongHourlyPrice: 10,
				LongHourlyStart: 300,
				LongHourlyEnd:   600,
			},
			wantErr: assert.NoError,
			mockExpect: func() {
				tariffRepoMock.AddMock.
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			uc := tariff.NewTariffUseCase(contextMock)

			_, err := uc.CreateTariff(tt.args)
			tt.wantErr(t, err)
		})
	}
}
