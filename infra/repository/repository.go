package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/internal/sql/queries"
)

type repositoryBase struct {
	db *sql.DB
}

func (r *repositoryBase) NewTransaction(ctx context.Context) (gateway.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

type accountRepository struct {
	repositoryBase
	queries *queries.Queries
}

func NewAccountRepository(db *sql.DB) gateway.AccountRepository {
	return &accountRepository{
		repositoryBase: repositoryBase{
			db: db,
		},
		queries: queries.New(db),
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account entity.Account) error {
	return r.queries.InsertNewAccount(ctx, queries.InsertNewAccountParams{
		ID:               account.ID,
		CustomerName:     account.CustomerName,
		DocumentNumber:   account.DocumentNUmber,
		Email:            account.Email,
		PasswordEncoded:  account.PasswordEncoded,
		SaltHashPassword: account.Salt,
		Status:           string(account.Status),
		AccountType:      string(account.AccountType),
		PhoneNumber:      account.PhoneNumber,
		CreatedAt:        account.CreatedAt,
		UpdatedAt:        account.UpdadatedAt,
	})
}

func (r *accountRepository) FindAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	row, err := r.queries.FindAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	account := r.rowToEntity(row.Account)
	raw := row.Transactions.(string)
	if err := json.Unmarshal(json.RawMessage(raw), &account.Wallet); err != nil {
		return nil, fmt.Errorf("repositoryBuildEntity: error in build entity from database, %w", err)
	}

	return &account, nil
}

func (r *accountRepository) FindAccountByIDs(ctx context.Context, ids ...uuid.UUID) (map[uuid.UUID]*entity.Account, error) {
	chError := make(chan error)
	chAccount := make(chan *entity.Account)

	for _, id := range ids {
		go func(accountID uuid.UUID) {
			account, err := r.FindAccount(ctx, accountID)
			chError <- err
			chAccount <- account
		}(id)
	}

	result := make(map[uuid.UUID]*entity.Account)
	for range ids {
		if err := <-chError; err != nil {
			return nil, err
		}

		account := <-chAccount
		result[account.ID] = account
	}

	return result, nil
}

func (r *accountRepository) SaveAtomicTransactions(ctx context.Context, transactions ...entity.Transaction) error {
	ch := make(chan error)
	for _, transaction := range transactions {
		ch <- r.query(ctx).InsertNewTransaction(ctx, queries.InsertNewTransactionParams{
			ID:              transaction.ID,
			CorrelatedID:    transaction.CorrelatedID,
			AccountID:       transaction.AccountID,
			TransactionType: string(transaction.TransactionType),
			Timestamp:       transaction.Timestamp,
			Amount:          int64(transaction.Amount),
		})
	}

	for range transactions {
		if err := <-ch; err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) rowToEntity(row queries.Account) entity.Account {
	return entity.Account{
		ID:              row.ID.(uuid.UUID),
		AccountType:     entity.AccountType(row.AccountType),
		CustomerName:    row.CustomerName,
		DocumentNUmber:  row.DocumentNumber,
		Email:           row.Email,
		PasswordEncoded: row.PasswordEncoded,
		Salt:            row.SaltHashPassword,
		PhoneNumber:     row.PhoneNumber,
		Status:          entity.AccountStatus(row.Status),
		CreatedAt:       row.CreatedAt.(time.Time),
		UpdadatedAt:     row.UpdatedAt.(time.Time),
	}
}

func (r *accountRepository) query(ctx context.Context) queries.Querier {
	tx, ok := gateway.GetTransactionContext(ctx)
	if !ok {
		return r.queries
	}
	txSQL := tx.(*sql.Tx)
	return r.queries.WithTx(txSQL)
}
