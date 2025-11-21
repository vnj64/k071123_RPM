package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/parking_service/contracts/pkg/proto"
	"k071123/internal/services/user_service/services/config"
	"log"
	"sync"
	"time"
)

var (
	parkingClient     proto.ParkingClient
	parkingClientOnce sync.Once
)

func makeSecureParkingConnection() (*grpc.ClientConn, error) {
	cfg := config.Make()
	addr := fmt.Sprintf("%s:%s", cfg.ParkingGrpcHost(), cfg.ParkingGrpcPort())
	// Используем контекст с тайм-аутом (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakeParkingServiceClient() (proto.ParkingClient, error) {
	var err error
	parkingClientOnce.Do(func() {
		conn, err := makeSecureParkingConnection()
		if err != nil {
			log.Fatalf("could not create client connection: %v", err)
		}
		parkingClient = proto.NewParkingClient(conn)
		log.Println("Parking client initialized successfully.")
	})
	return parkingClient, err
}
