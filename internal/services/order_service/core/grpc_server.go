package core

import (
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	grpc2 "k071123/internal/services/order_service/delivery/grpc"
	"k071123/internal/services/order_service/domain/cases/card"
	"k071123/internal/services/order_service/services/config"
	"k071123/internal/services/user_service/core"
	"net"
)

type GrpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	cfg := config.Make()
	ctx := InitCtx()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.OrderGrpcHost(), cfg.OrderGrpcPort()))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	nc, err := core.MakeNotificationServiceClient()
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	cardUseCase := card.NewCardUseCase(ctx, nc)
	proto.RegisterOrderServer(grpcServer, grpc2.NewHandler(ctx, cardUseCase))

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
