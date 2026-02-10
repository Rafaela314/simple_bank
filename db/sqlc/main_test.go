package db

import (
	"context"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"simple_bank/util"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stretchr/testify/require"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {

	var err error

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	exitCode := m.Run()

	if testDB != nil {
		testDB.Close()
	}

	os.Exit(exitCode)

}

// runTestWithTransaction runs a test function within a transaction that gets rolled back
func runTestWithTransaction(t *testing.T, testFunc func(*testing.T, *Queries)) {
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)

	txQueries := New(tx)

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Logf("Warning: failed to rollback transaction: %v", err)
		}
	}()

	// Run the test with transaction queries
	testFunc(t, txQueries)
}

// createTestData creates sample data for testing
func createTestData(t *testing.T) (Account, Entry, Transfer) {
	user := createRandomUser(t)
	currency := util.RandomCurrency()

	account := createRandomAccount(t, user.Username, currency)

	entry := createRandomEntryWithAccount(t, account)

	transfer := createRandomTransferWithAccounts(t, account, account)

	return account, entry, transfer
}

// createRandomAccount creates a random account for testing
func createRandomAccount(t *testing.T, owner string, currency string) Account {
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: currency,
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// createRandomEntryWithAccount creates an entry using the provided account
func createRandomEntryWithAccount(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

// createRandomTransferWithAccounts creates a transfer between accounts
func createRandomTransferWithAccounts(t *testing.T, fromAccount, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
