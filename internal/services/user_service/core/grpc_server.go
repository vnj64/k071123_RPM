package core

import (
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/user_service/contracts/pkg/proto"
	grpc2 "k071123/internal/services/user_service/delivery/grpc"
	"k071123/internal/services/user_service/domain/cases/user"
	"k071123/internal/services/user_service/services/config"
	"net"
)

type GrpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	cfg := config.Make()
	ctx := InitCtx()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.UserGrpcHost(), cfg.UserGrpcPort()))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	senderUseCase := user.NewUserUseCase(ctx)
	proto.RegisterUserServer(grpcServer, grpc2.NewHandler(ctx, senderUseCase))

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
