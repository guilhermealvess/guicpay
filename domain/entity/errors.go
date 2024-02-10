package entity

import (
	"fmt"

	"github.com/google/uuid"
)

type TransactionError struct {
	Message         string
	AccountID       uuid.UUID
	TransactionType TransactionType
	Amount          Money
}

func (t TransactionError) Error() string {
	return fmt.Sprintf("transaction_error: %s -> account_id=%s, transaction_type=%s, amount=%s", t.Message, t.AccountID, t.TransactionType, t.Amount)
}

func NewDepositError(msg string, accountID uuid.UUID, amount Money) *TransactionError {
	return &TransactionError{
		Message:         msg,
		AccountID:       accountID,
		TransactionType: Deposit,
		Amount:          amount,
	}
}

func NewTransferError(msg string, accountID uuid.UUID, amount Money) TransactionError {
	return TransactionError{
		Message:         msg,
		TransactionType: TransferPayer,
		AccountID:       accountID,
		Amount:          amount,
	}
}
