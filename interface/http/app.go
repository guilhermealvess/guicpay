package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewServer(h *accountHandler) *echo.Echo {
	otel.SetTextMapPropagator(trace.Baggage{})
	otel.SetTracerProvider(trace.NewNoopTracerProvider())
	server := echo.New()

	server.Use(otelecho.Middleware("my-server"))

	api := server.Group("/api")
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG")
	})

	api.POST("/accounts", h.CreateAccount)
	api.GET("/accounts", h.List)
	api.GET("/accounts/:id", h.Fetch)
	api.POST("/accounts/:id/deposit", h.AccountDeposit)
	api.POST("/accounts/:id/transfer", h.AccountTransfer)
	return server
}
