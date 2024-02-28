package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	"github.com/guilhermealvess/guicpay/internal/token"
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

func (h *accountHandler) Auth(c echo.Context) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	account, err := h.usecase.ExecuteLogin(c.Request().Context(), data.Email, data.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	raw := account.JsonRawMessage()
	t, err := token.JWT.Generate(raw)
	return buildResponse(c, err, t, http.StatusCreated)
}

func (h *accountHandler) CreateAccount(c echo.Context) error {
	var input usecase.NewAccountInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := usecase.ValidateDTO(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteNewAccount(c.Request().Context(), input)
	m := map[string]string{
		"account_id": output.String(),
	}

	return buildResponse(c, err, m, http.StatusCreated)
}

func (h *accountHandler) AccountDeposit(c echo.Context) error {
	v := c.Get(PayloadToken).(*Payload)
	var data struct {
		Value float64 `json:"value" validate:"required,min=0.01"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := usecase.ValidateDTO(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteDeposit(c.Request().Context(), v.AccountID, uint64(data.Value*100))
	m := map[string]string{
		"transaction_id": output.String(),
	}

	return buildResponse(c, err, m, http.StatusOK)
}

func (h *accountHandler) AccountTransfer(c echo.Context) error {
	v := c.Get(PayloadToken).(*Payload)
	var data struct {
		Value   float64   `json:"value" validate:"required,min=0.01"`
		PayeeID uuid.UUID `json:"payee" validate:"required"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := usecase.ValidateDTO(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	output, err := h.usecase.ExecuteTransfer(c.Request().Context(), v.AccountID, data.PayeeID, uint64(data.Value*100))
	m := map[string]string{
		"transaction_id": output.String(),
	}
	return buildResponse(c, err, m, http.StatusOK)
}

func (h *accountHandler) Fetch(c echo.Context) error {
	v := c.Get(PayloadToken).(*Payload)
	output, err := h.usecase.FindByID(c.Request().Context(), v.AccountID)
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
