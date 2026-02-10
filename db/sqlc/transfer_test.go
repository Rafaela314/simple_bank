package db

import (
	"context"
	"testing"

	"simple_bank/util"

	"github.com/stretchr/testify/require"
)

// TestCreateTransfer tests the creation of a transfer.
func TestCreateTransfer(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		currency := util.RandomCurrency()
		owner1 := createRandomUser(t).Username
		owner2 := createRandomUser(t).Username
		fromAccount := createAccountInTx(t, q, owner1, currency)
		toAccount := createAccountInTx(t, q, owner2, currency)
		arg := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        util.RandomMoney(),
		}

		transfer, err := q.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)

		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
		require.Equal(t, arg.Amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
	})
}

// TestGetTransfer tests the retrieval of a transfer.
func TestGetTransfer(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createTransferInTx(t, q)

		got, err := q.GetTransfer(context.Background(), created.ID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		require.Equal(t, created.ID, got.ID)
		require.Equal(t, created.FromAccountID, got.FromAccountID)
		require.Equal(t, created.ToAccountID, got.ToAccountID)
		require.Equal(t, created.Amount, got.Amount)
		require.Equal(t, created.CreatedAt, got.CreatedAt)
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		_, err := q.GetTransfer(context.Background(), 0)
		require.Error(t, err)
	})
}

// TestListTransfers tests the listing of transfers.
func TestListTransfers(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		currency := util.RandomCurrency()
		owner1 := createRandomUser(t).Username
		owner2 := createRandomUser(t).Username
		fromAccount := createAccountInTx(t, q, owner1, currency)
		toAccount := createAccountInTx(t, q, owner2, currency)
		n := 5
		var created []Transfer
		for i := 0; i < n; i++ {
			tr := createTransferInTxBetween(t, q, fromAccount.ID, toAccount.ID)
			created = append(created, tr)
		}

		listed, err := q.Listtransfers(context.Background(), ListtransfersParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Limit:         10,
			Offset:        0,
		})
		require.NoError(t, err)
		require.Len(t, listed, n)
		for i, tr := range listed {
			require.Equal(t, created[i].ID, tr.ID)
			require.True(t, tr.FromAccountID == fromAccount.ID || tr.ToAccountID == toAccount.ID)
		}
	})

	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		currency := util.RandomCurrency()
		owner1 := createRandomUser(t).Username
		owner2 := createRandomUser(t).Username
		fromAccount := createAccountInTx(t, q, owner1, currency)
		toAccount := createAccountInTx(t, q, owner2, currency)

		listed, err := q.Listtransfers(context.Background(), ListtransfersParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Limit:         10,
			Offset:        0,
		})
		require.NoError(t, err)
		require.Empty(t, listed)
	})
}

// TestUpdateTransfer tests the update of a transfer.
func TestUpdateTransfer(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createTransferInTx(t, q)
		newAmount := util.RandomMoney()

		updated, err := q.UpdateTransfer(context.Background(), UpdateTransferParams{
			ID:            created.ID,
			Amount:        newAmount,
			FromAccountID: created.FromAccountID,
			ToAccountID:   created.ToAccountID,
		})
		require.NoError(t, err)
		require.Equal(t, created.ID, updated.ID)
		require.Equal(t, newAmount, updated.Amount)
		require.Equal(t, created.FromAccountID, updated.FromAccountID)
		require.Equal(t, created.ToAccountID, updated.ToAccountID)

		got, err := q.GetTransfer(context.Background(), created.ID)
		require.NoError(t, err)
		require.Equal(t, newAmount, got.Amount)
	})
}

// TestDeleteTransfer tests the deletion of a transfer.
func TestDeletetransfers(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, q *Queries) {
		created := createTransferInTx(t, q)

		err := q.Deletetransfers(context.Background(), created.ID)
		require.NoError(t, err)

		_, err = q.GetTransfer(context.Background(), created.ID)
		require.Error(t, err)
	})
}

// createTransferInTx creates a transfer with random accounts.
func createTransferInTx(t *testing.T, q *Queries) Transfer {
	user1 := createRandomUser(t)
	owner1 := user1.Username
	user2 := createRandomUser(t)
	owner2 := user2.Username
	currency := util.RandomCurrency()
	fromAccount := createAccountInTx(t, q, owner1, util.RandomCurrency())
	toAccount := createAccountInTx(t, q, owner2, currency)
	return createTransferInTxBetween(t, q, fromAccount.ID, toAccount.ID)
}

// createTransferInTxBetween creates a transfer between two given accounts.
func createTransferInTxBetween(t *testing.T, q *Queries, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := q.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	return transfer
}
