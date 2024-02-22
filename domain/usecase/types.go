package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type NewAccountInput struct {
	Name           string `json:"customer_name" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Password       string `json:"password" validate:"required"`
	Type           string `json:"account_type" validate:"required"`
	DocumentNumber string `json:"document_number" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
}

type AccountOutput struct {
	ID           uuid.UUID `json:"account_id"`
	AccountType  string    `json:"account_type"`
	CustomerName string    `json:"customer_name"`
	Email        string    `json:"email"`
	Balance      string    `json:"balance"`
	Status       string    `json:"status"`
}

func ValidateDTO(v any) error {
	return validator.New().Struct(v)
}
