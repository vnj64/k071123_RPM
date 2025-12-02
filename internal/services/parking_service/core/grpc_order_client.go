package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/parking_service/services/config"
	"log"
	"sync"
	"time"
)

var (
	orderClient     proto.OrderClient
	orderClientOnce sync.Once
)

func makeSecureOrderConnection() (*grpc.ClientConn, error) {
	cfg := config.Make()
	addr := fmt.Sprintf("%s:%s", cfg.OrderGrpcHost(), cfg.OrderGrpcPort())
	// Используем контекст с тайм-аутом (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakeOrderServiceClient() (proto.OrderClient, error) {
	var err error
	orderClientOnce.Do(func() {
		conn, err := makeSecureOrderConnection()
		if err != nil {
			log.Fatalf("could not create client connection: %v", err)
		}
		orderClient = proto.NewOrderClient(conn)
		log.Println("Order client initialized successfully.")
	})
	return orderClient, err
}
