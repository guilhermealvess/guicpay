package gateway

import (
	"context"

	"github.com/guilhermealvess/guicpay/domain/entity"
)

type NotificationService interface {
	Notify(ctx context.Context, account entity.Account, transaction entity.Transaction) error
}

type AuthorizationService interface {
	Authorize(ctx context.Context, account entity.Account) error
}
