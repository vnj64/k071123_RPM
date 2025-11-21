package core

import (
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	grpc2 "k071123/internal/services/parking_service/delivery/grpc"
	"k071123/internal/services/parking_service/domain/cases/car"
	"k071123/internal/services/parking_service/services/config"
	"net"
)

type GrpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	cfg := config.Make()
	ctx := InitCtx()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.ParkingGrpcHost(), cfg.ParkingGrpcPort()))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	carUseCase := car.NewCarUseCase(ctx)
	proto.RegisterParkingServer(grpcServer, grpc2.NewHandler(ctx, carUseCase))

	return &GrpcServer{
		listener:   grpcListener,
		grpcServer: grpcServer,
	}
}

func (gs *GrpcServer) Start() {
	if err := gs.grpcServer.Serve(gs.listener); err != nil {
		panic("failed to serve: " + err.Error())
	}
}
