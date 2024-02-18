package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/internal/properties"
)

func (u *accountUseCase) ExecuteDeposit(ctx context.Context, accountID uuid.UUID, value uint64) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, properties.Props.TransactionTimeout)
	defer cancel()

	account, err := u.repository.FindAccount(ctx, accountID)
	if err != nil {
		return uuid.Nil, err
	}

	transaction, err := account.Deposit(entity.Money(value))
	if err != nil {
		return uuid.Nil, err
	}

	tx, err := u.repository.NewTransaction(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	ctx = gateway.InjectTransaction(ctx, tx)
	if err := u.repository.SaveAtomicTransactions(ctx, *transaction); err != nil {
		return uuid.Nil, err
	}

	if err := u.notification.Notify(ctx, *account, *transaction); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return transaction.ID, nil
}
