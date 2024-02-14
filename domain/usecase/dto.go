package usecase

import "github.com/guilhermealvess/guicpay/domain/entity"

type NewAccountInput struct {
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	Type           entity.AccountType `json:"account_type"`
	DocumentNumber string             `json:"document_number"`
	PhoneNumber    string             `json:"phone_number"`
}
