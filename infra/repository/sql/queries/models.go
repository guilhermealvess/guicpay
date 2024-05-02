package queries

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID              uuid.UUID `db:"id" json:"id"`
	AccountType     string    `db:"account_type" json:"account_type"`
	CustomerName    string    `db:"customer_name" json:"customer_name"`
	DocumentNumber  string    `db:"document_number" json:"document_number"`
	Email           string    `db:"email" json:"email"`
	PasswordEncoded string    `db:"password_encoded" json:"password_encoded"`
	PhoneNumber     string    `db:"phone_number" json:"phone_number"`
	Status          string    `db:"status" json:"status"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type Transaction struct {
	ID              uuid.UUID     `db:"id" json:"id"`
	AccountID       uuid.UUID     `db:"account_type" json:"account_type"`
	CorrelatedID    uuid.NullUUID `db:"correlated_id" json:"correlated_id"`
	Timestamp       time.Time     `db:"timestamp" json:"timestamp"`
	TransactionType string        `db:"transaction_type" json:"transaction_type"`
	Amount          int64         `db:"amount" json:"amount"`
	SnapshotID      uuid.NullUUID `db:"snapshot_id" json:"snapshot_id"`
	ParentID        uuid.UUID     `db:"parent_id" json:"parent_id"`
}

type ResumeAccount struct {
	ID          uuid.UUID `db:"id" json:"id"`
	AccountType string    `db:"account_type" json:"account_type"`
	Status      string    `db:"status" json:"status"`
	Email       string    `db:"email" json:"email"`
	Password    string    `db:"password_encoded" json:"password_encoded"`
}
