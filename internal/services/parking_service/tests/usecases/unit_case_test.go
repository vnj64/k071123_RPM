package usecases

import (
	"github.com/gojuno/minimock/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k071123/internal/services/parking_service/domain/cases/unit"
	mocks2 "k071123/internal/services/parking_service/domain/mocks"
	"k071123/internal/services/parking_service/domain/models/unit_statuses"
	"k071123/internal/services/parking_service/domain/props"
	"k071123/internal/services/parking_service/tests/mocks"
	"k071123/tools/logger"
	"testing"
)

func TestCarUseCase_CreateUnit(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	contextMock := mocks2.NewContextMock(mc)
	servicesMock := mocks2.NewServicesMock(mc)
	connMock := mocks2.NewConnectionMock(mc)

	unitRepoMock := mocks.NewUnitRepositoryMock(mc)

	realLogrus := logrus.New()
	fakeLog := &logger.Logger{
		Logger: realLogrus,
	}
	servicesMock.LoggerMock.Return(fakeLog)

	contextMock.ServicesMock.Return(servicesMock)
	contextMock.ConnectionMock.Return(connMock)

	connMock.UnitRepositoryMock.Return(unitRepoMock)

	// Success Create
	// Validation Error
	// DB Error
	tests := []struct {
		name       string
		args       props.CreateUniqReq
		wantErr    assert.ErrorAssertionFunc
		mockExpect func()
	}{
		{
			name: "Success Create (Required)",
			args: props.CreateUniqReq{
				Direction:     unit_statuses.In,
				NetworkStatus: unit_statuses.Online,
				Status:        unit_statuses.Active,
			},
			wantErr: assert.NoError,
			mockExpect: func() {
				unitRepoMock.AddMock.Return(nil)
			},
		},
		{
			name: "Validation Error",
			args: props.CreateUniqReq{
				Direction:     unit_statuses.In,
				NetworkStatus: unit_statuses.Online,
				Status:        unit_statuses.Active,
				Code:          ptrStr("12345"),
			},
			wantErr: assert.Error,
			mockExpect: func() {
				unitRepoMock.AddMock.Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			uc := unit.NewUnitUseCase(contextMock)

			_, err := uc.CreateUnit(tt.args)
			tt.wantErr(t, err)
		})
	}
}
