package fixture

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/guilhermealvess/guicpay/domain/common"
	"github.com/guilhermealvess/guicpay/domain/entity"
)

func FactoryFakeAccount() entity.Account {
	return entity.Account{
		ID:              uuid.New(),
		CustomerName:    gofakeit.Name(),
		DocumentNumber:  randomDocumentNumber(),
		Email:           gofakeit.Email(),
		AccountType:     entity.Personal,
		PhoneNumber:     gofakeit.Phone(),
		Status:          entity.AccountStatusActive,
		Wallet:          entity.Wallet{},
		CreatedAt:       gofakeit.Date(),
		UpdatedAt:       gofakeit.Date(),
		PasswordEncoded: entity.Password(fmt.Sprintf("SHA256:%s:%s", common.ComputeSHA256Hash("SALT"), common.ComputeSHA256Hash("TEST@PASSWORD"))),
	}
}

func FactoryFakeAccountWithBalance(balance entity.Money) entity.Account {
	account := FactoryFakeAccount()
	t := FactoryFakeTransaction(account, balance)
	account.Wallet = entity.Wallet{
		&t,
	}

	return account
}

func FactoryFakeTransaction(account entity.Account, amount entity.Money) entity.Transaction {
	return entity.Transaction{
		ID:              uuid.New(),
		AccountID:       account.ID,
		Amount:          amount.Absolute(),
		CorrelatedID:    uuid.NullUUID{},
		TransactionType: entity.Deposit,
		Timestamp:       gofakeit.Date(),
		SnapshotID:      uuid.NullUUID{},
	}
}

func FactoryFakeTransferOutput(payer, payee entity.Account, amount entity.Money) entity.TransferOutput {
	time := gofakeit.Date()
	t1 := entity.Transaction{
		ID:              uuid.New(),
		AccountID:       payer.ID,
		Amount:          amount.Absolute() * -1,
		CorrelatedID:    uuid.NullUUID{Valid: true, UUID: payer.ID},
		TransactionType: entity.TransferPayer,
		Timestamp:       time,
		SnapshotID:      uuid.NullUUID{},
	}

	t2 := entity.Transaction{
		ID:              uuid.New(),
		AccountID:       payee.ID,
		Amount:          amount.Absolute(),
		CorrelatedID:    uuid.NullUUID{Valid: true, UUID: payer.ID},
		TransactionType: entity.TransferPayee,
		Timestamp:       time,
		SnapshotID:      uuid.NullUUID{},
	}

	return entity.TransferOutput{
		Payer:        &t1,
		Payee:        &t2,
		CorrelatedID: t1.CorrelatedID.UUID,
	}
}

func randomDocumentNumber() string {
	return fmt.Sprintf("%d", gofakeit.Number(10000000000, 99999999999))
}
