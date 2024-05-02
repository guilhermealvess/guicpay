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
	ParentID        uuid.NullUUID
}

func factoryDepositTransaction(account Account, v Money) Transaction {
	transactionID := uuid.New()
	return Transaction{
		ID:              transactionID,
		AccountID:       account.ID,
		TransactionType: Deposit,
		Timestamp:       time.Now().UTC(),
		Amount:          v.Absolute(),
		ParentID:        uuid.NullUUID{Valid: true, UUID: transactionID},
	}
}

func factoryTransferTransactions(payerAccount, payeeAccount Account, v Money, parent *Transaction) (payer Transaction, payee Transaction) {
	now := time.Now().UTC()
	correlatedID := uuid.New()
	payer = Transaction{
		ID:              uuid.New(),
		CorrelatedID:    uuid.NullUUID{UUID: correlatedID, Valid: true},
		AccountID:       payerAccount.ID,
		TransactionType: TransferPayer,
		Timestamp:       now,
		Amount:          -1 * v.Absolute(),
		ParentID:        uuid.NullUUID{},
	}

	if parent != nil {
		payer.ParentID = uuid.NullUUID{Valid: true, UUID: parent.ID}
	}

	transactionPayeeID := uuid.New()
	payee = Transaction{
		ID:              transactionPayeeID,
		CorrelatedID:    uuid.NullUUID{UUID: correlatedID, Valid: true},
		AccountID:       payeeAccount.ID,
		TransactionType: TransferPayee,
		Timestamp:       now,
		Amount:          v.Absolute(),
		ParentID:        uuid.NullUUID{Valid: true, UUID: transactionPayeeID},
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

	t := &Transaction{
		ID:              snapshotID,
		AccountID:       accountID,
		TransactionType: Snapshot,
		Timestamp:       time.Now().UTC(),
		Amount:          balance,
		SnapshotID:      uuid.NullUUID{},
		CorrelatedID:    uuid.NullUUID{},
		ParentID:        uuid.NullUUID{},
	}

	if parent := w.FindParent(); parent != nil {
		t.ParentID = uuid.NullUUID{Valid: true, UUID: parent.ID}
	}

	return t
}

func (w *Wallet) FindParent() *Transaction {
	transactions := make([]Transaction, 0)
	m := make(map[uuid.UUID]bool)

	for _, t := range *w {
		if t.Amount < 0 {
			transactions = append(transactions, *t)
			m[t.ParentID.UUID] = true
		}
	}

	for _, t := range transactions {
		if _, ok := m[t.ID]; !ok {
			return &t
		}
	}

	return nil
}
