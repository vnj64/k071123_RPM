package main

import (
	"context"
	"k071123/internal/services/notification_service/core"
	"k071123/internal/services/notification_service/delivery/http"
	"k071123/internal/services/notification_service/services/amqp"
	"k071123/internal/services/notification_service/workers"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di := core.NewDi()
	log := di.Ctx.Services().Logger()
	server := core.NewHttpServer()
	grpcServer := core.NewGrpcServer()

	handlers := &http.Handlers{
		EmailSenderHandler: di.EmailSenderHandler,
	}

	workerManager := workers.NewWorkerManager()
	emailConsumer := amqp.NewEmailConsumer(di.Ctx.Services().Amqp().Publisher(), di.Ctx.Services().Smtp())
	workerManager.Register(emailConsumer)

	http.SetupRoutes(server.App(), handlers)

	wg.Add(1)
	go func() {
		if err := workerManager.StartAll(ctx); err != nil {
			log.Fatalf("failed to start workers: %v", err)
		}
		defer wg.Done()
		log.Infoln("Workers started")
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Infoln("HTTP server starting...")
		server.Start()
	}()
	go func() {
		defer wg.Done()
		log.Infoln("gRPC server starting...")
		grpcServer.Start()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Infoln("Shutting down notification service...")
	cancel()
	time.Sleep(time.Second)

	wg.Wait()
	log.Infoln("Notification service stopped gracefully")
}
