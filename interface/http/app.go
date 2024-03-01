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
	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG\n")
	})

	api := server.Group("/api")
	api.POST("/accounts", h.CreateAccount)
	api.GET("/accounts", h.List, validateToken)
	api.GET("/accounts/me", h.Fetch, validateToken)

	api.POST("/transactions/deposit", h.AccountDeposit, validateToken)
	api.POST("/transactions/transfer", h.AccountTransfer, validateToken)

	api.POST("/auth", h.Auth)

	return server
}
