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

func TestExecuteDeposit(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.Background()
		repository = mr.NewAccountRepository(t)
		tx         = mr.NewTx(t)
		usecase    = accountUseCase{
			repository: repository,
		}
		account = fixture.FactoryFakeAccount()
	)

	t.Run("Deposit", func(t *testing.T) {
		repository.On("FindAccount", ctx, account.ID).Return(&account, nil)
		repository.On("UpdateAccount", ctx, account).Return(nil)
		repository.On("NewTransaction", ctx).Return(tx, nil)
		tx.On("Commit").Return(nil)
		repository.On("SaveAtomicTransactions", ctx, mock.AnythingOfType("entity.Transaction")).Return(nil)

		id, err := usecase.ExecuteDeposit(ctx, account.ID, 100)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
	})
}
