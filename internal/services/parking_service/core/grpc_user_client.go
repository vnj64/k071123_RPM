package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/parking_service/services/config"
	"k071123/internal/services/user_service/contracts/pkg/proto"
	"log"
	"sync"
	"time"
)

var (
	userClient     proto.UserClient
	userClientOnce sync.Once
)

func makeSecureUserConnection() (*grpc.ClientConn, error) {
	cfg := config.Make()
	addr := fmt.Sprintf("%s:%s", cfg.UserGrpcHost(), cfg.UserGrpcPort())
	// Используем контекст с тайм-аутом (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakeUserServiceClient() (proto.UserClient, error) {
	var err error
	userClientOnce.Do(func() {
		conn, err := makeSecureUserConnection()
		if err != nil {
			log.Fatalf("could not create client connection: %v", err)
		}
		userClient = proto.NewUserClient(conn)
		log.Println("User client initialized successfully.")
	})
	return userClient, err
}
