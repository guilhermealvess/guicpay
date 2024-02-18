package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/labstack/echo/v4"
)

type accountHandler struct {
	usecase usecase.AccountUseCase
}

func NewAccountHandler(u usecase.AccountUseCase) *accountHandler {
	return &accountHandler{
		usecase: u,
	}
}

func (h *accountHandler) CreateAccount(c echo.Context) error {
	var input usecase.NewAccountInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteNewAccount(c.Request().Context(), input)
	m := map[string]string{
		"account_id": output.String(),
	}

	return buildResponse(c, err, m, http.StatusCreated)
}

func (h *accountHandler) AccountDeposit(c echo.Context) error {
	payeeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var data struct {
		Value float64 `json:"value"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteDeposit(c.Request().Context(), payeeID, uint64(data.Value*100))
	m := map[string]string{
		"transaction_id": output.String(),
	}

	return buildResponse(c, err, m, http.StatusOK)
}

func (h *accountHandler) AccountTransfer(c echo.Context) error {
	payerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var data struct {
		Value   float64   `json:"value"`
		PayeeID uuid.UUID `json:"payee"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteTransfer(c.Request().Context(), payerID, data.PayeeID, uint64(data.Value*100))
	m := map[string]string{
		"transaction_id": output.String(),
	}
	return buildResponse(c, err, m, http.StatusOK)
}

func (h *accountHandler) Fetch(c echo.Context) error {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.FindByID(c.Request().Context(), accountID)
	return buildResponse(c, err, output, http.StatusOK)
}

func (h *accountHandler) List(c echo.Context) error {
	output, err := h.usecase.FindAll(c.Request().Context())
	return buildResponse(c, err, output, http.StatusOK)
}

func buildResponse(c echo.Context, err error, data any, statusCode int) error {
	switch {
	case err == nil:
		return c.JSON(statusCode, data)

	case errors.Is(err, sql.ErrNoRows):
		return echo.NewHTTPError(http.StatusBadRequest, err)

	case errors.Is(err, entity.ErrUnprocessableEntity):
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
}
