package http

import (
	"net/http"

	_ "github.com/guilhermealvess/guicpay/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewServer(h *accountHandler) *echo.Echo {
	server := echo.New()
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	api := server.Group("/api")
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG\n")
	})

	api.POST("/accounts", h.CreateAccount)
	api.GET("/accounts", h.List)
	api.GET("/accounts/:id", h.Fetch)
	api.POST("/accounts/:id/deposit", h.AccountDeposit)
	api.POST("/accounts/:id/transfer", h.AccountTransfer)
	return server
}
