package repository

import (
	"context"
	"testing"
	"time"

	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/fixture"
	"github.com/guilhermealvess/guicpay/infra/repository/sql/queries"
	"github.com/guilhermealvess/guicpay/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRepository(t *testing.T) {
	var (
		db         = getDatabaseContainerConnection(t)
		repository = NewAccountRepository(db)
		ctx        = context.Background()
	)

	t.Run("CreateAccount", func(t *testing.T) {
		defer resetDatabase(t, db)
		account := fixture.FactoryFakeAccount()

		err := repository.CreateAccount(context.Background(), account)
		assert.NoError(t, err)

		const query = `SELECT * FROM accounts WHERE id = $1`
		var accountDB queries.Account
		err = db.Get(&accountDB, query, account.ID)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, accountDB.ID)
		assert.Equal(t, account.CustomerName, accountDB.CustomerName)
		assert.Equal(t, account.DocumentNumber, accountDB.DocumentNumber)
		assert.Equal(t, account.Email, accountDB.Email)
		assert.Equal(t, account.PhoneNumber, accountDB.PhoneNumber)
		assert.Equal(t, string(account.Status), accountDB.Status)
		assert.Equal(t, string(account.AccountType), accountDB.AccountType)
		assert.Equal(t, account.CreatedAt.Format(time.DateTime), accountDB.CreatedAt.Format(time.DateTime))
		assert.Equal(t, account.UpdatedAt.Format(time.DateTime), accountDB.UpdatedAt.Format(time.DateTime))
		assert.Equal(t, string(account.PasswordEncoded), accountDB.PasswordEncoded)
	})

	t.Run("FindAccount", func(t *testing.T) {
		defer resetDatabase(t, db)
		account := fixture.FactoryFakeAccountWithBalance(entity.MilReais)
		err := repository.CreateAccount(context.Background(), account)
		assert.NoError(t, err)

		accountDB, err := repository.FindAccount(context.Background(), account.ID)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, accountDB.ID)
		assert.Equal(t, account.CustomerName, accountDB.CustomerName)
		assert.Equal(t, account.DocumentNumber, accountDB.DocumentNumber)
		assert.Equal(t, account.Email, accountDB.Email)
		assert.Equal(t, account.PhoneNumber, accountDB.PhoneNumber)
		assert.Equal(t, account.Status, accountDB.Status)
		assert.Equal(t, account.AccountType, accountDB.AccountType)
		assert.Equal(t, account.PasswordEncoded, accountDB.PasswordEncoded)
	})

	t.Run("FindAccountByIDs", func(t *testing.T) {
		defer resetDatabase(t, db)
		account1 := fixture.FactoryFakeAccount()
		account2 := fixture.FactoryFakeAccount()
		err := repository.CreateAccount(ctx, account1)
		assert.NoError(t, err)
		err = repository.CreateAccount(ctx, account2)
		assert.NoError(t, err)

		accounts, err := repository.FindAccountByIDs(ctx, account1.ID, account2.ID)
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)
		assert.NotNil(t, accounts[account1.ID])
		assert.NotNil(t, accounts[account2.ID])
	})

	t.Run("SaveAtomicTransactions", func(t *testing.T) {
		defer resetDatabase(t, db)
		account1 := fixture.FactoryFakeAccountWithBalance(entity.MilReais)
		account2 := fixture.FactoryFakeAccountWithBalance(entity.MilReais / 2)
		accounts := []entity.Account{account1, account2}

		err := repository.SaveAtomicTransactions(ctx, *account1.Wallet[0], *account2.Wallet[0])

		assert.NoError(t, err)

		for _, account := range accounts {
			const query = `SELECT * FROM transaction WHERE account_id = $1`
			var transaction queries.Transaction
			err = db.Get(&transaction, query, account.ID)
			assert.NoError(t, err)
			assert.Equal(t, account.ID, transaction.AccountID)
			assert.Equal(t, account.Wallet[0].Amount, transaction.Amount)
			assert.Equal(t, account.Wallet[0].CorrelatedID, transaction.CorrelatedID)
			assert.Equal(t, account.Wallet[0].TransactionType, transaction.TransactionType)
			assert.Equal(t, account.Wallet[0].Timestamp, transaction.Timestamp)
			assert.Equal(t, account.Wallet[0].SnapshotID, transaction.SnapshotID)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		defer resetDatabase(t, db)
		account1 := fixture.FactoryFakeAccount()
		account2 := fixture.FactoryFakeAccount()
		err := repository.CreateAccount(ctx, account1)
		assert.NoError(t, err)
		err = repository.CreateAccount(ctx, account2)
		assert.NoError(t, err)

		accounts, err := repository.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)
	})

	t.Run("FindAccountByEmail", func(t *testing.T) {
		defer resetDatabase(t, db)
		account := fixture.FactoryFakeAccount()
		err := repository.CreateAccount(ctx, account)
		assert.NoError(t, err)

		account2 := fixture.FactoryFakeAccount()
		err = repository.CreateAccount(ctx, account2)
		assert.NoError(t, err)

		accountDB, err := repository.FindAccountByEmail(ctx, account.Email)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, accountDB.ID)
		assert.Equal(t, account.CustomerName, accountDB.CustomerName)
		assert.Equal(t, account.DocumentNumber, accountDB.DocumentNumber)
		assert.Equal(t, account.Email, accountDB.Email)
		assert.Equal(t, account.PhoneNumber, accountDB.PhoneNumber)
		assert.Equal(t, account.Status, accountDB.Status)
		assert.Equal(t, account.AccountType, accountDB.AccountType)
		assert.Equal(t, account.CreatedAt, accountDB.CreatedAt)
		assert.Equal(t, account.UpdatedAt, accountDB.UpdatedAt)
		assert.Equal(t, account.PasswordEncoded, accountDB.PasswordEncoded)
		assert.NotEqual(t, accountDB.Email, account2.Email)
	})
}

func depositMoney(t *testing.T, repository accountRepository, account entity.Account, value entity.Money) {
	t.Helper()
	transaction, err := account.Deposit(value)
	assert.NoError(t, err)

	err = repository.SaveAtomicTransactions(context.Background(), *transaction)
	assert.NoError(t, err)
}

func getDatabaseContainerConnection(t testing.TB) *sqlx.DB {
	t.Helper()
	pgContainer, err := postgres.RunContainer(
		context.Background(),
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts("./sql/schema.sql"),
	)

	if err != nil {
		t.Fatalf("could not start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		t.Fatalf("could not get connection string: %v", err)
	}

	return database.NewConnectionDB(connStr)
}

func resetDatabase(t testing.TB, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		t.Fatalf("could not delete from transaction: %v", err)
	}

	_, err = db.Exec("DELETE FROM accounts")
	if err != nil {
		t.Fatalf("could not delete from account: %v", err)
	}
}
