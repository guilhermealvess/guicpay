package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/fixture"
	mg "github.com/guilhermealvess/guicpay/mocks/domain/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteTransfer(t *testing.T) {
	var (
		ctx          = context.Background()
		repository   = mg.NewAccountRepository(t)
		tx           = mg.NewTx(t)
		mutex        = mg.NewMutex(t)
		auth         = mg.NewAuthorizationService(t)
		notification = mg.NewNotificationService(t)
		usecase      = accountUseCase{
			repository:   repository,
			authorizer:   auth,
			notification: notification,
			mutex:        mutex,
			queue:        make(chan uuid.UUID),
		}
	)

	setupMock := func() {
		repository.Mock = mock.Mock{}
		mutex.Mock = mock.Mock{}
		auth.Mock = mock.Mock{}
		notification.Mock = mock.Mock{}
		tx.Mock = mock.Mock{}
	}

	t.Run("Trasnfer with successful", func(t *testing.T) {
		t.Parallel()
		setupMock()

		var (
			accountFrom            = fixture.FactoryFakeAccountWithBalance(entity.MilReais)
			accountTo              = fixture.FactoryFakeAccount()
			FindAccountByIDsResult = map[uuid.UUID]*entity.Account{
				accountFrom.ID: &accountFrom,
				accountTo.ID:   &accountTo,
			}
		)

		repository.On("NewTransaction", mock.Anything).Return(tx, nil)
		repository.On("FindAccountByIDs", mock.Anything, accountFrom.ID, accountTo.ID).Return(FindAccountByIDsResult, nil)
		repository.On("SaveAtomicTransactions", mock.Anything, mock.AnythingOfType("entity.Transaction"), mock.AnythingOfType("entity.Transaction")).Return(nil)
		tx.On("Commit").Return(nil)

		mutex.On("Lock", mock.Anything, accountFrom.ID.String(), mock.AnythingOfType("time.Duration")).Return(nil)
		mutex.On("Unlock", mock.Anything, accountFrom.ID.String()).Return(nil)

		notification.On("Notify", mock.Anything, mock.MatchedBy(func(a entity.Account) bool {
			return a.ID == accountTo.ID
		}), mock.AnythingOfType("entity.Transaction")).Return(nil)
		auth.On("Authorize", mock.Anything, accountFrom).Return(nil)

		id, err := usecase.ExecuteTransfer(ctx, accountFrom.ID, accountTo.ID, 100)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
		tx.AssertNotCalled(t, "Rollback")

		assert.Len(t, usecase.queue, 0)
	})

	t.Run("Trasnfer with error", func(t *testing.T) {
		setupMock()
		var (
			accountFrom = fixture.FactoryFakeAccountWithBalance(entity.MilReais)
			accountTo   = fixture.FactoryFakeAccount()
		)

		mutex.On("Lock", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
		repository.On("NewTransaction", mock.Anything).Return(tx, nil)
		tx.On("Commit").Return(nil)

		id, err := usecase.ExecuteTransfer(ctx, accountFrom.ID, accountTo.ID, 100)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		tx.AssertNumberOfCalls(t, "Commit", 1)
		repository.AssertNotCalled(t, "FindAccountByIDs")
		repository.AssertNotCalled(t, "SaveAtomicTransactions")
		mutex.AssertNotCalled(t, "Unlock")
		notification.AssertNotCalled(t, "Notify")
		auth.AssertNotCalled(t, "Authorize")

		assert.Len(t, usecase.queue, 0)
	})
}
