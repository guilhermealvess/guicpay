package usecase

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/guilhermealvess/guicpay/domain/fixture"
	mr "github.com/guilhermealvess/guicpay/mocks/domain/gateway"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAuthentication(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.Background()
		repository = mr.NewAccountRepository(t)
		usecase    = accountUseCase{
			repository: repository,
		}

		email   = gofakeit.Email()
		account = fixture.FactoryFakeAccount()
	)

	repository.On("FindResumeAccount", ctx, email).Return(&account, nil)
	output, err := usecase.ExecuteLogin(ctx, email, "123456")

	assert.NoError(t, err)
	assert.Equal(t, account, *output)
}
