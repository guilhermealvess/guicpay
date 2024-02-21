package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
)

type AccountRepository interface {
	Repository
	CreateAccount(ctx context.Context, account entity.Account) error
	FindAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error)
	FindAccountByIDs(ctx context.Context, ids ...uuid.UUID) (map[uuid.UUID]*entity.Account, error)
	SaveAtomicTransactions(ctx context.Context, transactions ...entity.Transaction) error
	FindAll(ctx context.Context) ([]*entity.Account, error)
	SetSnapshotTransactions(ctx context.Context, snapshotID uuid.UUID, transactionIDs uuid.UUIDs) error
}

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository interface {
	NewTransaction(ctx context.Context) (Tx, error)
}

type transactionContextKey string

const TransactionContextKey transactionContextKey = "TransactionContextKey"

func InjectTransaction(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, TransactionContextKey, tx)
}

func GetTransactionContext(ctx context.Context) (Tx, bool) {
	val := ctx.Value(TransactionContextKey)
	if val == nil {
		return nil, false
	}

	tx := val.(Tx)
	return tx, true
}
