package grpcport

import (
	"context"

	"github.com/google/uuid"
)

type AccountContext string

const AccountContextKey AccountContext = "AccountContext"

func getAccountContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(AccountContextKey)
	if v == nil {
		return uuid.Nil, false
	}

	return v.(uuid.UUID), true
}

func setAccountContext(ctx context.Context, accountID uuid.UUID) context.Context {
	return context.WithValue(ctx, AccountContextKey, accountID)
}
