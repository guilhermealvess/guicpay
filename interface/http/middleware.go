package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/internal/token"
	"github.com/labstack/echo/v4"
)

const PayloadToken = "account"

type Payload struct {
	AccountID   uuid.UUID `json:"account_id"`
	AccountType string    `json:"account_type"`
}

func validateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		var payload Payload
		if err := token.Middle(tokenString, &payload); err != nil {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		c.Set(PayloadToken, &payload)
		return next(c)
	}
}

func validatePermission(accountType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload := c.Get(PayloadToken).(*Payload)
			if payload.AccountType != accountType {
				return echo.NewHTTPError(http.StatusForbidden)
			}

			return next(c)
		}
	}
}
