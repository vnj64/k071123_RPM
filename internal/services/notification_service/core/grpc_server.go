package core

import (
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/notification_service/contracts/pkg/proto"
	grpc2 "k071123/internal/services/notification_service/delivery/grpc"
	"k071123/internal/services/notification_service/domain/cases/sender"
	"k071123/internal/services/notification_service/services/config"
	"net"
)

type GrpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

// TODO: прийти к синглтону domain.Context
func NewGrpcServer() *GrpcServer {
	cfg := config.Make()
	ctx := InitCtx()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.NotificationGrpcHost(), cfg.NotificationGrpcPort()))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	senderUseCase := sender.NewEmailSenderUseCase(ctx)
	proto.RegisterNotificationServer(grpcServer, grpc2.NewHandler(ctx, senderUseCase))

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
