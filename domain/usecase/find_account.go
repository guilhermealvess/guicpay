package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (u *accountUseCase) FindByID(ctx context.Context, accountID uuid.UUID) (*AccountOutput, error) {
	account, err := u.repository.FindAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &AccountOutput{
		ID:           account.ID,
		AccountType:  string(account.AccountType),
		CustomerName: account.CustomerName,
		Email:        account.Email,
		Status:       string(account.Status),
		Balance:      account.Wallet.Balance().String(),
	}, nil
}

func (u *accountUseCase) FindAll(ctx context.Context) ([]*AccountOutput, error) {
	accounts, err := u.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*AccountOutput, 0)
	for _, account := range accounts {
		data := AccountOutput{
			ID:           account.ID,
			AccountType:  string(account.AccountType),
			CustomerName: account.CustomerName,
			Email:        account.Email,
			Status:       string(account.Status),
			Balance:      account.Wallet.Balance().String(),
		}
		result = append(result, &data)
	}

	return result, nil
}
