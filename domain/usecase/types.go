package usecase

import (
	"github.com/google/uuid"
)

type NewAccountInput struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Type           string `json:"account_type"`
	DocumentNumber string `json:"document_number"`
	PhoneNumber    string `json:"phone_number"`
}

type AccountOutput struct {
	ID           uuid.UUID `json:"account_id"`
	AccountType  string    `json:"account_type"`
	CustomerName string    `json:"customer_name"`
	Email        string    `json:"email"`
	Balance      string    `json:"balance"`
	Status       string    `json:"status"`
}
