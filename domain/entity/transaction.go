package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Deposit       TransactionType = "DEPOSIT"
	TransferPayer TransactionType = "TRANSFER_PAYER"
	TransferPayee TransactionType = "TRANSFER_PAYEE"
	Snapshot      TransactionType = "SNAPSHOT"
)

type Transaction struct {
	ID              uuid.UUID
	CorrelatedID    uuid.NullUUID
	AccountID       uuid.UUID
	TransactionType TransactionType
	Timestamp       time.Time
	Amount          Money
	SnapshotID      uuid.NullUUID
}

func factoryDepositTransaction(account Account, v Money) Transaction {
	t := Transaction{
		ID:              uuid.New(),
		AccountID:       account.ID,
		TransactionType: Deposit,
		Timestamp:       time.Now().UTC(),
		Amount:          v.Absolute(),
	}

	return t
}

func factoryTransferTransactions(payerAccount, payeeAccount Account, v Money) (payer Transaction, payee Transaction) {
	now := time.Now().UTC()
	correlatedID := uuid.New()
	payer = Transaction{
		ID:              uuid.New(),
		CorrelatedID:    uuid.NullUUID{UUID: correlatedID, Valid: true},
		AccountID:       payerAccount.ID,
		TransactionType: TransferPayer,
		Timestamp:       now,
		Amount:          -1 * v.Absolute(),
	}

	payee = Transaction{
		ID:              uuid.New(),
		CorrelatedID:    uuid.NullUUID{UUID: correlatedID, Valid: true},
		AccountID:       payeeAccount.ID,
		TransactionType: TransferPayee,
		Timestamp:       now,
		Amount:          v.Absolute(),
	}

	return
}

type Wallet []*Transaction

func (w *Wallet) Balance() Money {
	var balance Money
	for _, transaction := range *w {
		balance += transaction.Amount
	}

	return balance
}

func (w *Wallet) Snapshot(accountID uuid.UUID) *Transaction {
	snapshotID := uuid.New()
	var balance Money

	for _, t := range *w {
		balance += t.Amount
		t.SnapshotID = uuid.NullUUID{UUID: snapshotID, Valid: true}
	}

	return &Transaction{
		ID:              snapshotID,
		AccountID:       accountID,
		TransactionType: Snapshot,
		Timestamp:       time.Now().UTC(),
		Amount:          balance,
		SnapshotID:      uuid.NullUUID{},
		CorrelatedID:    uuid.NullUUID{},
	}
}
