package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	const (
		name          = "Fulano De Tal"
		cpf           = "1235678910"
		cnpj          = "00623904000173"
		emailPersonal = "personal@example.com"
		emailSeller   = "seller@example.com"
		pass          = "PASSWORD"
		phone         = "+5511996344108"
	)

	personal := NewAccount(
		Personal,
		name,
		cpf,
		emailPersonal,
		pass,
		phone,
	)

	seller := NewAccount(
		Seller,
		name,
		cnpj,
		emailSeller,
		pass,
		phone,
	)

	assert.NotEqual(t, uuid.Nil, personal.ID)
	assert.NotEqual(t, uuid.Nil, seller.ID)

	t.Run("deposit", func(t *testing.T) {
		account := Account(personal)
		v := 300*Real + 55*Cent
		before := time.Now()
		tr, err := account.Deposit(v)

		assert.NoError(t, err)
		assert.Equal(t, v, tr.Amount)
		assert.NotEqual(t, uuid.Nil, tr.ID)
		assert.False(t, tr.CorrelatedID.Valid)
		assert.True(t, before.Before(tr.Timestamp))
		assert.Equal(t, account.ID, tr.AccountID)
		assert.Equal(t, Deposit, tr.TransactionType)
		assert.Len(t, account.Wallet, 1)
		assert.Equal(t, account.Wallet[0].ID, tr.ID)
	})

	t.Run("failure deposit", func(t *testing.T) {
		account := Account(personal)
		account.Status = AccountStatusCanceled
		v := 300*Real + 55*Cent

		tr, err := account.Deposit(v)
		assert.Error(t, err)
		assert.Nil(t, tr)
	})

	t.Run("transfer", func(t *testing.T) {
		pa := Account(personal)
		v := 300*Real + 55*Cent
		depositInAccount(&pa, v)

		before := time.Now()
		sa := Account(seller)
		output, err := pa.Transfer(&sa, v)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, output.CorrelatedID)
		assert.True(t, output.CorrelatedID == output.Payee.CorrelatedID.UUID && output.CorrelatedID == output.Payer.CorrelatedID.UUID)
		assert.NotEqual(t, uuid.Nil, output.Payer.ID)
		assert.Equal(t, output.Payer.AccountID, pa.ID)
		assert.True(t, before.Before(output.Payer.Timestamp))
		assert.True(t, -1*v == output.Payer.Amount)
		assert.True(t, v == output.Payee.Amount)
	})
}

func depositInAccount(account *Account, v Money) {
	t := factoryDepositTransaction(*account, v)
	account.Wallet = append(account.Wallet, t)
}
