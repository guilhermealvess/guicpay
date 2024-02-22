package queries

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Query interface {
	GetContext(context.Context, any, string, ...any) error
	SelectContext(context.Context, any, string, ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Queries struct {
	db Query
}

func New(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx *sqlx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

type FindAccountByIDRow struct {
	Account
	Transactions json.RawMessage `db:"transactions"`
}

func (q *Queries) FindAccountByID(ctx context.Context, id uuid.UUID) (*FindAccountByIDRow, error) {
	const findAccountByID = `SELECT ac.id, 
		ac.account_type, 
		ac.customer_name, 
		ac.document_number, 
		ac.email, 
		ac.password_encoded, 
		ac.salt_hash_password, 
		ac.phone_number, 
		ac.status, 
		ac.created_at, 
		ac.updated_at,
		CASE
			WHEN tr.account_id IS NULL THEN 'null'::json
			ELSE json_agg(tr.*)
		END AS transactions
	FROM accounts ac
	LEFT JOIN transactions tr ON ac.id = tr.account_id
	WHERE ac.id = $1 AND tr.snapshot_id IS NULL GROUP BY ac.id, tr.account_id;`
	var row FindAccountByIDRow
	if err := q.db.GetContext(ctx, &row, findAccountByID, id); err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	return &row, nil
}

type SaveAccountParams struct {
	ID              uuid.UUID `db:"id" json:"id"`
	AccountType     string    `db:"account_type" json:"account_type"`
	CustomerName    string    `db:"customer_name" json:"customer_name"`
	DocumentNumber  string    `db:"document_number" json:"document_number"`
	Email           string    `db:"email" json:"email"`
	PasswordEncoded string    `db:"password_encoded" json:"password_encoded"`
	SaltHash        string    `db:"salt_hash_password" json:"salt_hash_password"`
	PhoneNumber     string    `db:"phone_number" json:"phone_number"`
	Status          string    `db:"status" json:"status"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

func (q *Queries) SaveAccount(ctx context.Context, params SaveAccountParams) error {
	const query = `
	INSERT INTO accounts (id,account_type,customer_name,document_number,email,password_encoded,salt_hash_password,phone_number,status,created_at,updated_at) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);`
	_, err := q.db.ExecContext(ctx, query, params.ID, params.AccountType, params.CustomerName, params.DocumentNumber, params.Email, params.PasswordEncoded, params.SaltHash, params.PhoneNumber, params.Status, params.CreatedAt, params.UpdatedAt)
	return err
}

type SaveTransactionParams struct {
	ID              uuid.UUID     `db:"id" json:"id"`
	CorrelatedID    uuid.NullUUID `db:"correlated_id" json:"correlated_id"`
	AccountID       uuid.UUID     `db:"account_id" json:"account_id"`
	TransactionType string        `db:"transaction_type" json:"transaction_type"`
	Timestamp       time.Time     `db:"timestamp" json:"timestamp"`
	Amount          int64         `db:"amount" json:"amount"`
	SnapshotID      uuid.NullUUID `db:"snapshot_id" json:"snapshot_id"`
}

func (q *Queries) SaveTransaction(ctx context.Context, params SaveTransactionParams) error {
	const query = `INSERT INTO transactions (id,correlated_id,account_id,transaction_type,timestamp,amount,snapshot_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := q.db.ExecContext(ctx, query, params.ID, params.CorrelatedID, params.AccountID, params.TransactionType, params.Timestamp, params.Amount, params.SnapshotID)
	return err
}

func (q *Queries) FindAll(ctx context.Context) ([]*FindAccountByIDRow, error) {
	const query = `SELECT ac.id, 
		ac.account_type, 
		ac.customer_name, 
		ac.document_number, 
		ac.email, 
		ac.password_encoded, 
		ac.salt_hash_password, 
		ac.phone_number, 
		ac.status, 
		ac.created_at, 
		ac.updated_at,
		CASE
			WHEN tr.account_id IS NULL THEN 'null'::json
			ELSE json_agg(tr.*)
		END AS transactions
	FROM accounts ac LEFT JOIN transactions tr on tr.account_id = ac.id 
	WHERE tr.snapshot_id IS NULL
	GROUP BY ac.id, tr.account_id ORDER BY ac.created_at desc;`

	var rows []*FindAccountByIDRow
	if err := q.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	return rows, nil
}

func (q *Queries) SetSnapshotTransactions(ctx context.Context, snapshotID uuid.UUID, transactionIDs uuid.UUIDs) error {
	query := "UPDATE transactions SET snapshot_id = $1 WHERE id IN ("
	for i, id := range transactionIDs {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("'%s'", id.String())
	}
	query += ")"

	result, err := q.db.ExecContext(ctx, query, snapshotID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("database: %w", sql.ErrNoRows)
	}

	return nil
}
