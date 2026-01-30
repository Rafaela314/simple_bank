package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		arg := CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}

		account, err := q.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, account)

		require.Equal(t, arg.Owner, account.Owner)
		require.Equal(t, arg.Balance, account.Balance)
		require.Equal(t, arg.Currency, account.Currency)

		require.NotZero(t, account.ID)
		require.NotZero(t, account.CreatedAt)
	})
}

func TestGetAccount(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)

		got, err := q.GetAccount(context.Background(), created.ID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		require.Equal(t, created.ID, got.ID)
		require.Equal(t, created.Owner, got.Owner)
		require.Equal(t, created.Balance, got.Balance)
		require.Equal(t, created.Currency, got.Currency)
		require.Equal(t, created.CreatedAt, got.CreatedAt)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		_, err := q.GetAccount(context.Background(), 0)
		require.Error(t, err)
	})
}

func TestGetAccountForUpdate(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)

		got, err := q.GetAccountForUpdate(context.Background(), created.ID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		require.Equal(t, created.ID, got.ID)
		require.Equal(t, created.Owner, got.Owner)
		require.Equal(t, created.Balance, got.Balance)
		require.Equal(t, created.Currency, got.Currency)
	})
}

func TestListAccounts(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		owner := util.RandomOwner()
		n := 5
		var created []Account
		for i := 0; i < n; i++ {
			acc := createAccountInTxWithOwner(t, q, owner)
			created = append(created, acc)
		}

		listed, err := q.ListAccounts(context.Background(), ListAccountsParams{
			Owner:  owner,
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err)
		require.Len(t, listed, n)
		for i, acc := range listed {
			require.Equal(t, created[i].ID, acc.ID)
			require.Equal(t, owner, acc.Owner)
		}
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		listed, err := q.ListAccounts(context.Background(), ListAccountsParams{
			Owner:  "nonexistent_owner_xyz",
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err)
		require.Empty(t, listed)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		owner := util.RandomOwner()
		for i := 0; i < 5; i++ {
			createAccountInTxWithOwner(t, q, owner)
		}

		listed, err := q.ListAccounts(context.Background(), ListAccountsParams{
			Owner:  owner,
			Limit:  2,
			Offset: 1,
		})
		require.NoError(t, err)
		require.Len(t, listed, 2)
	})
}

func TestUpdateAccount(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)
		newBalance := util.RandomMoney()

		updated, err := q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      created.ID,
			Balance: newBalance,
		})
		require.NoError(t, err)
		require.Equal(t, created.ID, updated.ID)
		require.Equal(t, created.Owner, updated.Owner)
		require.Equal(t, newBalance, updated.Balance)
		require.Equal(t, created.Currency, updated.Currency)

		got, err := q.GetAccount(context.Background(), created.ID)
		require.NoError(t, err)
		require.Equal(t, newBalance, got.Balance)
	})
}

func TestAddAccountBalance(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)
		amount := int64(50)

		updated, err := q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     created.ID,
			Amount: amount,
		})
		require.NoError(t, err)
		require.Equal(t, created.ID, updated.ID)
		require.Equal(t, created.Balance+amount, updated.Balance)

		got, err := q.GetAccount(context.Background(), created.ID)
		require.NoError(t, err)
		require.Equal(t, created.Balance+amount, got.Balance)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)
		withdraw := int64(-30)

		updated, err := q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     created.ID,
			Amount: withdraw,
		})
		require.NoError(t, err)
		require.Equal(t, created.Balance+withdraw, updated.Balance)
	})
}

func TestDeleteAccount(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createAccountInTx(t, q)

		err := q.DeleteAccount(context.Background(), created.ID)
		require.NoError(t, err)

		_, err = q.GetAccount(context.Background(), created.ID)
		require.Error(t, err)
	})
}

// createAccountInTx creates a random account using the given Queries (e.g. transaction-scoped).
func createAccountInTx(t *testing.T, q *Queries) Account {
	return createAccountInTxWithOwner(t, q, util.RandomOwner())
}

func createAccountInTxWithOwner(t *testing.T, q *Queries, owner string) Account {
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := q.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	return account
}
