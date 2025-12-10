package session

import (
	notifyOrder "k071123/internal/services/notification_service/contracts/pkg/proto"
	orderProto "k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/props"
	userProto "k071123/internal/services/user_service/contracts/pkg/proto"
	"reflect"
	"testing"
)

func TestSessionUseCase_Start(t *testing.T) {
	type fields struct {
		ctx domain.Context
		oc  orderProto.OrderClient
		nc  notifyOrder.NotificationClient
		uc  userProto.UserClient
	}
	type args struct {
		args props.StartSessionReq
		oc   orderProto.OrderClient
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp props.StartSessionResp
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &SessionUseCase{
				ctx: tt.fields.ctx,
				oc:  tt.fields.oc,
				nc:  tt.fields.nc,
				uc:  tt.fields.uc,
			}
			gotResp, err := uc.Start(tt.args.args, tt.args.oc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Start() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
