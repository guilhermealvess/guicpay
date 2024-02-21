package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/internal/logger"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"go.uber.org/zap"
)

func (u *accountUseCase) ExecuteTransfer(ctx context.Context, payer, payee uuid.UUID, value uint64) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, properties.Props.TransactionTimeout)
	defer cancel()

	tx, err := u.repository.NewTransaction(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	ctx = gateway.InjectTransaction(ctx, tx)
	if err := u.mutex.Lock(ctx, payer.String(), properties.Props.TransactionTimeout); err != nil {
		return uuid.Nil, err
	}

	defer func() {
		if err := u.mutex.Unlock(ctx, payer.String()); err != nil {
			tx.Rollback()
			logger.Logger.Warn("Error in unlock payer account", zap.String("account_id", payer.String()))
		}
	}()

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
