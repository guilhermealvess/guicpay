package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"go.opentelemetry.io/otel"
)

func (u *accountUseCase) ExecuteTransfer(ctx context.Context, payer, payee uuid.UUID, value uint64) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, properties.Props.TransactionTimeout)
	defer cancel()

	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "AccountUseCase.ExecuteTransfer")
	defer span.End()

	tx, err := u.repository.NewTransaction(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	defer tx.Commit()
	ctx = gateway.InjectTransaction(ctx, tx)
	accounts, err := u.repository.FindAccountByIDs(ctx, payer, payee)
	if err != nil {
		return uuid.Nil, err
	}

	payerAccount, payeeAccount := accounts[payer], accounts[payee]
	if err := u.authorizer.Authorize(ctx, *payerAccount); err != nil {
		return uuid.Nil, err
	}

	output, err := payerAccount.Transfer(payeeAccount, entity.Money(value))
	if err != nil {
		return uuid.Nil, err
	}

	if len(payerAccount.Wallet) >= properties.Props.SnapshotWalletSize {
		go func() {
			u.queue <- payerAccount.ID
		}()
	}

	if err := u.repository.SaveAtomicTransactions(ctx, *output.Payer, *output.Payee); err != nil {
		return uuid.Nil, err
	}

	if err := u.notification.Notify(ctx, *payeeAccount, *output.Payee); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return output.CorrelatedID, tx.Commit()
}
