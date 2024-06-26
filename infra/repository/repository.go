package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/infra/repository/sql/queries"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type repositoryBase struct {
	db *sqlx.DB
}

func (r *repositoryBase) NewTransaction(ctx context.Context) (gateway.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

type accountRepository struct {
	repositoryBase
	queries *queries.Queries
}

func NewAccountRepository(db *sqlx.DB) gateway.AccountRepository {
	return &accountRepository{
		repositoryBase: repositoryBase{
			db: db,
		},
		queries: queries.New(db),
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account entity.Account) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "CreateAccount")
	defer span.End()
	err := r.queries.SaveAccount(ctx, queries.SaveAccountParams{
		ID:              account.ID,
		CustomerName:    account.CustomerName,
		DocumentNumber:  account.DocumentNumber,
		Email:           account.Email,
		PasswordEncoded: string(account.PasswordEncoded),
		Status:          string(account.Status),
		AccountType:     string(account.AccountType),
		PhoneNumber:     account.PhoneNumber,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdatedAt,
	})

	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (r *accountRepository) FindAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "FindAccount")
	defer span.End()

	row, err := r.queries.FindAccountByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("database: %w", err)
	}

	account := entity.Account{
		ID:              row.Account.ID,
		AccountType:     entity.AccountType(row.Account.AccountType),
		CustomerName:    row.Account.CustomerName,
		DocumentNumber:  row.Account.DocumentNumber,
		Email:           row.Account.Email,
		PasswordEncoded: entity.Password(row.Account.PasswordEncoded),
		PhoneNumber:     row.Account.PhoneNumber,
		Status:          entity.AccountStatus(row.Account.Status),
		CreatedAt:       row.Account.CreatedAt,
		UpdatedAt:       row.Account.UpdatedAt,
	}

	if err := json.Unmarshal(row.Transactions, &account.Wallet); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("database: %w", err)
	}

	return &account, nil
}

func (r *accountRepository) FindAccountByIDs(ctx context.Context, ids ...uuid.UUID) (map[uuid.UUID]*entity.Account, error) {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "FindAccountByIDs")
	defer span.End()

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
			span.RecordError(err)
			return nil, fmt.Errorf("database: %w", err)
		}

		account := <-chAccount
		result[account.ID] = account
	}

	return result, nil
}

func (r *accountRepository) SaveAtomicTransactions(ctx context.Context, transactions ...entity.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "SaveAtomicTransactions")
	defer span.End()
	ch := make(chan error)
	for _, t := range transactions {
		go func(transaction entity.Transaction) {
			ch <- r.query(ctx).SaveTransaction(ctx, queries.SaveTransactionParams{
				ID:              transaction.ID,
				CorrelatedID:    transaction.CorrelatedID,
				AccountID:       transaction.AccountID,
				TransactionType: string(transaction.TransactionType),
				Timestamp:       transaction.Timestamp,
				Amount:          int64(transaction.Amount),
				ParentID:        transaction.ParentID,
			})
		}(t)
	}

	for range transactions {
		if err := <-ch; err != nil {
			span.RecordError(err)
			return fmt.Errorf("database: %w", err)
		}
	}

	return nil
}

func (r *accountRepository) FindAll(ctx context.Context) ([]*entity.Account, error) {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "FindAll")
	defer span.End()

	rows, err := r.queries.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	accounts := make([]*entity.Account, 0)
	for _, row := range rows {
		account := entity.Account{
			ID:              row.Account.ID,
			AccountType:     entity.AccountType(row.Account.AccountType),
			CustomerName:    row.Account.CustomerName,
			DocumentNumber:  row.Account.DocumentNumber,
			Email:           row.Account.Email,
			PasswordEncoded: entity.Password(row.Account.PasswordEncoded),
			PhoneNumber:     row.Account.PhoneNumber,
			Status:          entity.AccountStatus(row.Account.Status),
			CreatedAt:       row.Account.CreatedAt,
			UpdatedAt:       row.Account.UpdatedAt,
		}

		if err := json.Unmarshal(row.Transactions, &account.Wallet); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("database: %w", err)
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (r *accountRepository) SetSnapshotTransactions(ctx context.Context, snapshotID uuid.UUID, transactionIDs uuid.UUIDs) error {
	return r.query(ctx).SetSnapshotTransactions(ctx, snapshotID, transactionIDs)
}

func (r *accountRepository) FindAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "FindAccountByEmail")
	defer span.End()

	row, err := r.queries.FindAccountByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	account := entity.Account{
		ID:              row.Account.ID,
		AccountType:     entity.AccountType(row.Account.AccountType),
		CustomerName:    row.Account.CustomerName,
		DocumentNumber:  row.Account.DocumentNumber,
		Email:           row.Account.Email,
		PasswordEncoded: entity.Password(row.Account.PasswordEncoded),
		PhoneNumber:     row.Account.PhoneNumber,
		Status:          entity.AccountStatus(row.Account.Status),
		CreatedAt:       row.Account.CreatedAt,
		UpdatedAt:       row.Account.UpdatedAt,
	}

	if err := json.Unmarshal(row.Transactions, &account.Wallet); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("database: %w", err)
	}

	return &account, nil
}

func (r *accountRepository) FindResumeAccount(ctx context.Context, email string) (*entity.ResumeAccount, error) {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "FindResumeAccount")
	defer span.End()

	row, err := r.queries.FindResumeAccount(ctx, email)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	account := entity.ResumeAccount{
		ID:              row.ID,
		AccountType:     entity.AccountType(row.AccountType),
		Email:           row.Email,
		Status:          entity.AccountStatus(row.Status),
		PasswordEncoded: entity.Password(row.Password),
	}

	return &account, nil
}

func (r *accountRepository) query(ctx context.Context) *queries.Queries {
	tx, ok := gateway.GetTransactionContext(ctx)
	if !ok {
		return r.queries
	}
	txSQL := tx.(*sqlx.Tx)
	return r.queries.WithTx(txSQL)
}
