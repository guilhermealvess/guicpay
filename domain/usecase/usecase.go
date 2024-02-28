package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
)

type AccountUseCase interface {
	ExecuteNewAccount(ctx context.Context, input NewAccountInput) (uuid.UUID, error)
	ExecuteDeposit(ctx context.Context, accountID uuid.UUID, value uint64) (uuid.UUID, error)
	ExecuteTransfer(ctx context.Context, payer, payee uuid.UUID, value uint64) (uuid.UUID, error)
	FindByID(ctx context.Context, accountID uuid.UUID) (*AccountOutput, error)
	FindAll(ctx context.Context) ([]*AccountOutput, error)
	ExecuteSnapshotTransaction(ctx context.Context, accountID uuid.UUID)
	ExecuteLogin(ctx context.Context, email, password string) (*entity.ResumeAccount, error)
}

type accountUseCase struct {
	repository   gateway.AccountRepository
	authorizer   gateway.AuthorizationService
	notification gateway.NotificationService
	mutex        gateway.Mutex
	queue        chan uuid.UUID
}

func NewAccountUseCase(r gateway.AccountRepository, m gateway.Mutex, n gateway.NotificationService, a gateway.AuthorizationService, ch chan uuid.UUID) AccountUseCase {
	return &accountUseCase{
		repository:   r,
		authorizer:   a,
		notification: n,
		mutex:        m,
		queue:        ch,
	}
}
