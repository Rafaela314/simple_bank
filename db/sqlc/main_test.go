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

	dns := os.Getenv("DB_SOURCE")
	if dns == "" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		db := os.Getenv("POSTGRES_DB")
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}

		if user == "" || password == "" || db == "" {
			log.Fatal("cannot connect to database. database configuration is missing; set DB_SOURCE or POSTGRES_* env vars")
		}
		dns = "postgresql://" + user + ":" + password + "@" + host + ":5432/" + db + "?sslmode=disable"
	}

	testDB, err = pgxpool.New(context.Background(), dns)
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

	account := createRandomAccount(t, user.Username, util.RandomCurrency())

	entry := createRandomEntryWithAccount(t, account)

	transfer := createRandomTransferWithAccounts(t, account, account)

	return account, entry, transfer
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
