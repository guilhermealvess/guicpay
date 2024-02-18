package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/infra/repository/sql/queries"
	"github.com/jmoiron/sqlx"
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
	return r.queries.SaveAccount(ctx, queries.SaveAccountParams{
		ID:              account.ID,
		CustomerName:    account.CustomerName,
		DocumentNumber:  account.DocumentNumber,
		Email:           account.Email,
		PasswordEncoded: account.PasswordEncoded,
		SaltHash:        account.Salt,
		Status:          string(account.Status),
		AccountType:     string(account.AccountType),
		PhoneNumber:     account.PhoneNumber,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdadatedAt,
	})
}

func (r *accountRepository) FindAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	row, err := r.queries.FindAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	createdAt, err := row.Account.CreatedAt.Time()
	if err != nil {
		return nil, err
	}

	updatedAt, err := row.Account.UpdatedAt.Time()
	if err != nil {
		return nil, err
	}

	account := entity.Account{
		ID:              row.Account.ID,
		AccountType:     entity.AccountType(row.Account.AccountType),
		CustomerName:    row.Account.CustomerName,
		DocumentNumber:  row.Account.DocumentNumber,
		Email:           row.Account.Email,
		PasswordEncoded: row.Account.PasswordEncoded,
		Salt:            row.Account.Salt,
		PhoneNumber:     row.Account.PhoneNumber,
		Status:          entity.AccountStatus(row.Account.Status),
		CreatedAt:       createdAt,
		UpdadatedAt:     updatedAt,
	}

	var transactions []queries.Transaction
	if err := row.Transactions.Bind(&transactions); err != nil {
		return nil, err
	}

	for _, t := range transactions {
		timestamp, err := t.Timestamp.Time()
		if err != nil {
			return nil, err
		}

		transaction := entity.Transaction{
			ID:              t.ID,
			CorrelatedID:    t.CorrelatedID,
			AccountID:       t.AccountID,
			TransactionType: entity.TransactionType(t.TransactionType),
			Timestamp:       timestamp,
			Amount:          entity.Money(t.Amount),
		}

		account.Wallet = append(account.Wallet, transaction)
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
	for _, t := range transactions {
		go func(transaction entity.Transaction) {
			ch <- r.query(ctx).SaveTransaction(ctx, queries.SaveTransactionParams{
				ID:              transaction.ID,
				CorrelatedID:    transaction.CorrelatedID,
				AccountID:       transaction.AccountID,
				TransactionType: string(transaction.TransactionType),
				Timestamp:       transaction.Timestamp,
				Amount:          int64(transaction.Amount),
			})
		}(t)
	}

	for range transactions {
		if err := <-ch; err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) FindAll(ctx context.Context) ([]*entity.Account, error) {
	rows, err := r.queries.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make([]*entity.Account, 0)
	for _, row := range rows {
		createdAt, err := row.Account.CreatedAt.Time()
		if err != nil {
			return nil, err
		}

		updatedAt, err := row.Account.CreatedAt.Time()
		if err != nil {
			return nil, err
		}

		account := entity.Account{
			ID:              row.Account.ID,
			AccountType:     entity.AccountType(row.Account.AccountType),
			CustomerName:    row.Account.CustomerName,
			DocumentNumber:  row.Account.DocumentNumber,
			Email:           row.Account.Email,
			PasswordEncoded: row.Account.PasswordEncoded,
			Salt:            row.Account.Salt,
			PhoneNumber:     row.Account.PhoneNumber,
			Status:          entity.AccountStatus(row.Account.Status),
			CreatedAt:       createdAt,
			UpdadatedAt:     updatedAt,
		}

		var transactions []queries.Transaction
		if err := row.Transactions.Bind(&transactions); err != nil {
			return nil, err
		}

		for _, t := range transactions {
			timestamp, err := t.Timestamp.Time()
			if err != nil {
				return nil, err
			}

			transaction := entity.Transaction{
				ID:              t.ID,
				CorrelatedID:    t.CorrelatedID,
				AccountID:       t.AccountID,
				TransactionType: entity.TransactionType(t.TransactionType),
				Timestamp:       timestamp,
				Amount:          entity.Money(t.Amount),
			}

			account.Wallet = append(account.Wallet, transaction)
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (r *accountRepository) query(ctx context.Context) *queries.Queries {
	tx, ok := gateway.GetTransactionContext(ctx)
	if !ok {
		return r.queries
	}
	txSQL := tx.(*sqlx.Tx)
	return r.queries.WithTx(txSQL)
}