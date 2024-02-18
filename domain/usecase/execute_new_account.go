package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
)

func (u *accountUseCase) ExecuteNewAccount(ctx context.Context, input NewAccountInput) (uuid.UUID, error) {
	account := entity.NewAccount(
		entity.AccountType(strings.ToUpper(input.Type)),
		input.Name,
		input.DocumentNumber,
		input.Email,
		input.Password,
		input.PhoneNumber,
	)

	if err := u.repository.CreateAccount(ctx, account); err != nil {
		return uuid.Nil, err
	}

	return account.ID, nil
}
