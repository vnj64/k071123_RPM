package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/order_service/services/config"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	"log"
	"sync"
	"time"
)

var (
	client     proto.ParkingClient
	clientOnce sync.Once
)

func makeSecureConnection() (*grpc.ClientConn, error) {
	cfg := config.Make()
	// Используем контекст с тайм-аутом (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := fmt.Sprintf("%s:%s", cfg.ParkingGrpcHost(), cfg.ParkingGrpcPort())
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakeParkingServiceClient() (proto.ParkingClient, error) {
	var err error
	clientOnce.Do(func() {
		conn, err := makeSecureConnection()
		if err != nil {
			log.Fatalf("could not create client connection: %v", err)
		}
		client = proto.NewParkingClient(conn)
		log.Println("Parking client initialized successfully.")
	})
	return client, err
}
