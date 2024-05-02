package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	name          = "Fulano De Tal"
	cpf           = "1235678910"
	cnpj          = "00623904000173"
	emailPersonal = "personal@example.com"
	emailSeller   = "seller@example.com"
	pass          = "PASSWORD"
	phone         = "+5511996344108"
)

func TestTransaction(t *testing.T) {
	seller := factoryFakeSellerAccount(t)
	personal := factoryFakePersonalAccount(t)
	value := 300 * Real

	t.Run("Deposit", func(t *testing.T) {
		transaction := factoryDepositTransaction(personal, value)

		assert.NotEqual(t, uuid.Nil, transaction.ID)
		assert.Equal(t, personal.ID, transaction.AccountID)
		assert.False(t, transaction.CorrelatedID.Valid)
		assert.Equal(t, Deposit, transaction.TransactionType)
		assert.Equal(t, value.Absolute(), transaction.Amount)
	})

	t.Run("Transfer", func(t *testing.T) {
		payerTransaction, payeeTransaction := factoryTransferTransactions(personal, seller, value, nil)

		assert.NotEqual(t, uuid.Nil, payerTransaction.CorrelatedID.UUID)
		assert.Equal(t, payerTransaction.CorrelatedID, payeeTransaction.CorrelatedID)
		assert.True(t, payerTransaction.Amount < 0 && payerTransaction.Amount.Absolute() == value)
		assert.Equal(t, payerTransaction.AccountID, personal.ID)

		assert.Equal(t, payeeTransaction.AccountID, seller.ID)
	})
}

func factoryFakePersonalAccount(t testing.TB) Account {
	t.Helper()
	personal := NewAccount(
		Personal,
		name,
		cpf,
		emailPersonal,
		pass,
		phone,
	)

	assert.NotEqual(t, uuid.Nil, personal.ID)
	return personal
}

func factoryFakeSellerAccount(t testing.TB) Account {
	t.Helper()
	seller := NewAccount(
		Seller,
		name,
		cnpj,
		emailSeller,
		pass,
		phone,
	)

	assert.NotEqual(t, uuid.Nil, seller.ID)
	return seller
}
