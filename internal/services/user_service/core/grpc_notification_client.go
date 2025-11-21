package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k071123/internal/services/notification_service/contracts/pkg/proto"
	"k071123/internal/services/user_service/services/config"
	"log"
	"sync"
	"time"
)

var (
	notificationClient     proto.NotificationClient
	notificationClientOnce sync.Once
)

func makeSecureNotificationConnection() (*grpc.ClientConn, error) {
	cfg := config.Make()
	addr := fmt.Sprintf("%s:%s", cfg.NotificationGrpcHost(), cfg.NotificationGrpcPort())
	// Используем контекст с тайм-аутом (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakeNotificationServiceClient() (proto.NotificationClient, error) {
	var err error
	notificationClientOnce.Do(func() {
		conn, err := makeSecureNotificationConnection()
		if err != nil {
			log.Fatalf("could not create client connection: %v", err)
		}
		notificationClient = proto.NewNotificationClient(conn)
		log.Println("Notification client initialized successfully.")
	})
	return notificationClient, err
}
