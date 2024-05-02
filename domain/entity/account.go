package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	Personal AccountType = "PERSONAL"
	Seller   AccountType = "SELLER"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "ACTIVE"
	AccountStatusCanceled AccountStatus = "CANCELED"
)

type Account struct {
	ID              uuid.UUID
	AccountType     AccountType
	CustomerName    string
	DocumentNumber  string
	Email           string
	PasswordEncoded Password
	PhoneNumber     string
	Status          AccountStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Wallet          Wallet
}

func NewAccount(t AccountType, name, doc, email, pass, phone string) Account {
	now := time.Now().UTC()
	return Account{
		ID:              uuid.New(),
		AccountType:     t,
		CustomerName:    name,
		DocumentNumber:  doc,
		Email:           email,
		PasswordEncoded: generatePasswordEncoded(pass),
		PhoneNumber:     phone,
		Status:          AccountStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
		Wallet:          []*Transaction{},
	}
}

func (a *Account) Deposit(v Money) (*Transaction, error) {
	if a.Status == AccountStatusCanceled {
		return nil, errors.Join(ErrUnprocessableEntity, NewDepositError("account canceled cant deposit", a.ID, v))
	}

	t := factoryDepositTransaction(*a, v)
	a.Wallet = append(a.Wallet, &t)

	return &t, nil
}

func (a *Account) Transfer(payee *Account, v Money) (*TransferOutput, error) {
	if a.AccountType == Seller {
		return nil, errors.Join(ErrUnprocessableEntity, NewTransferError("account seller cant make transfer", a.ID, v))
	}

	if a.ID == payee.ID {
		return nil, errors.Join(ErrUnprocessableEntity, NewTransferError("account cant transfer to itself", a.ID, v))
	}

	if a.Wallet.Balance() < v {
		return nil, errors.Join(ErrUnprocessableEntity, NewTransferError("insuficient balance", a.ID, v))
	}

	t1, t2 := factoryTransferTransactions(*a, *payee, v, a.Wallet.FindParent())
	a.Wallet = append(a.Wallet, &t1)
	payee.Wallet = append(payee.Wallet, &t2)

	return &TransferOutput{
		Payer:        &t1,
		Payee:        &t2,
		CorrelatedID: t1.CorrelatedID.UUID,
	}, nil
}

type TransferOutput struct {
	Payer        *Transaction
	Payee        *Transaction
	CorrelatedID uuid.UUID
}

type ResumeAccount struct {
	ID              uuid.UUID
	AccountType     AccountType
	Email           string
	Status          AccountStatus
	PasswordEncoded Password
}

func (a *ResumeAccount) ValidatePassword(pass string) error {
	return a.PasswordEncoded.Compare(pass)
}

func (a *ResumeAccount) JsonRawMessage() json.RawMessage {
	raw, _ := json.Marshal(map[string]string{"account_id": a.ID.String()})
	return raw
}
