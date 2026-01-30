package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		account := createAccountInTx(t, q)
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}

		entry, err := q.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)

		require.Equal(t, arg.AccountID, entry.AccountID)
		require.Equal(t, arg.Amount, entry.Amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	})
}

func TestGetEntry(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createEntryInTx(t, q)

		got, err := q.GetEntry(context.Background(), created.ID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		require.Equal(t, created.ID, got.ID)
		require.Equal(t, created.AccountID, got.AccountID)
		require.Equal(t, created.Amount, got.Amount)
		require.Equal(t, created.CreatedAt, got.CreatedAt)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		_, err := q.GetEntry(context.Background(), 0)
		require.Error(t, err)
	})
}

func TestListEntries(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		account := createAccountInTx(t, q)
		n := 5
		var created []Entry
		for i := 0; i < n; i++ {
			entry := createEntryInTxForAccount(t, q, account.ID)
			created = append(created, entry)
		}

		listed, err := q.ListEntries(context.Background(), ListEntriesParams{
			AccountID: account.ID,
			Limit:     10,
			Offset:    0,
		})
		require.NoError(t, err)
		require.Len(t, listed, n)
		for i, entry := range listed {
			require.Equal(t, created[i].ID, entry.ID)
			require.Equal(t, account.ID, entry.AccountID)
		}
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		account := createAccountInTx(t, q)

		listed, err := q.ListEntries(context.Background(), ListEntriesParams{
			AccountID: account.ID,
			Limit:     10,
			Offset:    0,
		})
		require.NoError(t, err)
		require.Empty(t, listed)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		account := createAccountInTx(t, q)
		for i := 0; i < 5; i++ {
			createEntryInTxForAccount(t, q, account.ID)
		}

		listed, err := q.ListEntries(context.Background(), ListEntriesParams{
			AccountID: account.ID,
			Limit:     2,
			Offset:    1,
		})
		require.NoError(t, err)
		require.Len(t, listed, 2)
	})
}

func createEntryInTx(t *testing.T, q *Queries) Entry {
	account := createAccountInTx(t, q)
	return createEntryInTxForAccount(t, q, account.ID)
}

func createEntryInTxForAccount(t *testing.T, q *Queries, accountID int64) Entry {
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}
	entry, err := q.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	return entry
}
