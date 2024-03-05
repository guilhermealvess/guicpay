package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/fixture"
	mr "github.com/guilhermealvess/guicpay/mocks/domain/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteNewAccount(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		repository = mr.NewAccountRepository(t)
		u          = accountUseCase{
			repository: repository,
		}
	)

	t.Run("CreateAccount", func(t *testing.T) {
		account := fixture.FactoryFakeAccount()
		repository.On("CreateAccount", ctx, mock.AnythingOfType("entity.Account")).Return(nil)

		id, err := u.ExecuteNewAccount(ctx, NewAccountInput{
			Name:           account.CustomerName,
			Email:          account.Email,
			Type:           string(account.AccountType),
			Password:       "teste123",
			DocumentNumber: account.DocumentNumber,
			PhoneNumber:    account.PhoneNumber,
		})

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
	})
}
