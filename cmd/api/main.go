package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/guilhermealvess/guicpay/infra/mutex"
	"github.com/guilhermealvess/guicpay/infra/repository"
	"github.com/guilhermealvess/guicpay/infra/service"
	grpcport "github.com/guilhermealvess/guicpay/interface/grpc_port"
	"github.com/guilhermealvess/guicpay/interface/http"
	"github.com/guilhermealvess/guicpay/internal/database"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Guic Pay Simplificado ...")

	queue, snapshotBackgroundWorker := buildSnapShotWorker()

	// Gateway
	repo := repository.NewAccountRepository(database.NewConnectionDB())
	mu := mutex.NewMutex(properties.Props.RedisAddress, "")
	notificationService := service.NewNotificationService(properties.Props.NotificationServiceURL)
	authService := service.NewAuthorizationService(properties.Props.AuthorizeServiceURL)

	// UseCase
	usecase := usecase.NewAccountUseCase(repo, mu, notificationService, authService, queue)
	go snapshotBackgroundWorker(usecase)

	// Handler
	handler := http.NewAccountHandler(usecase)

	// Application Server
	server := http.NewServer(handler)
	server.Use(middleware.Logger())
	server.Use(middleware.RequestID())
	server.Use(middleware.Recover())
	server.Use(middleware.AddTrailingSlash())
	go func() {
		server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", properties.Props.RestPort)))
	}()

	grpcApp := grpcport.NewApp(usecase)
	grpcApp.Start(properties.Props.GRPCPort)
	close(queue)
}

func buildSnapShotWorker() (chan uuid.UUID, func(usecase.AccountUseCase)) {
	queue := make(chan uuid.UUID)

	return queue, func(usecase usecase.AccountUseCase) {
		for accountID := range queue {
			go usecase.ExecuteSnapshotTransaction(context.Background(), accountID)
		}
	}
}
