package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
)

func (u *accountUseCase) ExecuteDeposit(ctx context.Context, accountID uuid.UUID, value uint64) (uuid.UUID, error) {
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

	if err := u.notification.Notify(ctx, *account, *transaction); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return transaction.ID, tx.Commit()
}
