package usecase

import (
	"context"

	"github.com/guilhermealvess/guicpay/domain/entity"
)

func (u *accountUseCase) ExecuteLogin(ctx context.Context, email, password string) (*entity.ResumeAccount, error) {
	account, err := u.repository.FindResumeAccount(ctx, email)
	if err != nil {
		return nil, err
	}

	return account, account.ValidatePassword(password)
}
