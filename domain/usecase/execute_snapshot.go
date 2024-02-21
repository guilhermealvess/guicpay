package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	"github.com/guilhermealvess/guicpay/internal/logger"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"go.uber.org/zap"
)

func (u *accountUseCase) ExecuteSnapshotTransaction(ctx context.Context, accountID uuid.UUID) {
	ctx, cancel := context.WithTimeout(ctx, properties.Props.TransactionTimeout)
	defer cancel()

	tx, err := u.repository.NewTransaction(ctx)
	if err != nil {
		logger.Logger.Error("Error in new tx", zap.Error(err))
		return
	}

	ctx = gateway.InjectTransaction(ctx, tx)
	if err := u.mutex.Lock(ctx, accountID.String(), properties.Props.TransactionTimeout); err != nil {
		logger.Logger.Error("Error in mutex lock", zap.Error(err))
		return
	}

	defer func() {
		if err := u.mutex.Unlock(ctx, accountID.String()); err != nil {
			logger.Logger.Error("Error in rollback", zap.Error(err))
			tx.Rollback()
		}
	}()

	account, err := u.repository.FindAccount(ctx, accountID)
	if err != nil {
		logger.Logger.Error("Error in find account", zap.Error(err))
		return
	}

	if len(account.Wallet) < properties.Props.SnapshotWalletSize {
		return
	}

	snapshot := account.Wallet.Snapshot(account.ID)
	if err := u.repository.SaveAtomicTransactions(ctx, *snapshot); err != nil {
		logger.Logger.Error("Error in save snapshot", zap.Error(err))
		return
	}

	transactionIDs := make(uuid.UUIDs, 0)
	for _, t := range account.Wallet {
		transactionIDs = append(transactionIDs, t.ID)
	}

	if err := u.repository.SetSnapshotTransactions(ctx, snapshot.ID, transactionIDs); err != nil {
		logger.Logger.Error("Error in update snapshot", zap.Error(err))
		tx.Rollback()
		return
	}

	tx.Commit()
	logger.Logger.Info("Done snapshot", zap.String("snapshotID", snapshot.ID.String()), zap.String("account_id", account.ID.String()))
}
