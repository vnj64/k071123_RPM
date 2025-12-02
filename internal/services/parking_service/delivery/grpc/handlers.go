package grpc

import (
	"context"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/cases/car"
	"k071123/internal/services/parking_service/domain/props"
)

type GrpcHandler struct {
	proto.UnimplementedParkingServer
	ctx        domain.Context
	carUseCase *car.CarUseCase
}

func NewHandler(ctx domain.Context, carUseCase *car.CarUseCase) *GrpcHandler {
	return &GrpcHandler{
		ctx:        ctx,
		carUseCase: carUseCase,
	}
}

func (h *GrpcHandler) CreateCar(c context.Context, req *proto.CreateCarReq) (*proto.CreateCarResp, error) {
	resp, err := h.carUseCase.CreateCar(props.CreateCarReq{
		GosNumber: req.GosNumber,
		UserUUID:  req.UserUUID,
	})
	if err != nil {
		return nil, err
	}
	return &proto.CreateCarResp{
		Car: &proto.Car{
			UUID:      resp.Car.UUID.String(),
			UserUUID:  resp.Car.UserUUID.String(),
			GosNumber: resp.Car.GosNumber,
			IsActive:  resp.Car.IsActive,
		},
	}, nil
}
