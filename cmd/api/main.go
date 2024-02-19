package main

import (
	"fmt"

	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/guilhermealvess/guicpay/infra/mutex"
	"github.com/guilhermealvess/guicpay/infra/repository"
	"github.com/guilhermealvess/guicpay/infra/service"
	"github.com/guilhermealvess/guicpay/interface/http"
	"github.com/guilhermealvess/guicpay/internal/database"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Guic Pay Simplificado ...")

	// Gateway
	repo := repository.NewAccountRepository(database.NewConnectionDB())
	mu := mutex.NewMutex(properties.Props.RedisAddress, "")
	notificationService := service.NewNotificationService(properties.Props.NotificationServiceURL)
	authService := service.NewAuthorizationService(properties.Props.AuthorizeServiceURL)

	// UseCase
	usecase := usecase.NewAccountUseCase(repo, mu, notificationService, authService)

	// Handler
	handler := http.NewAccountHandler(usecase)

	// Application Server
	server := http.NewServer(handler)
	server.Use(middleware.Logger())
	server.Use(middleware.RequestID())
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", properties.Props.Port)))
}
