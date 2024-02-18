package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewServer(h *accountHandler) *echo.Echo {
	server := echo.New()
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
