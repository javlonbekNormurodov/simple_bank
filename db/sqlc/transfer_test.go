package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createTransferTest(t *testing.T) Transfer {
	account1, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{Limit: 2, Offset: 2})
	require.NoError(t, err)
	require.Len(t, account1, 2)
	args := CreateTransferParams{
		FromAccountID: account1[0].ID,
		ToAccountID:   account1[1].ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Amount, args.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createTransferTest(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createTransferTest(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createTransferTest(t)
	args := UpdateTransferParams{
		ID:            transfer1.ID,
		FromAccountID: transfer1.FromAccountID,
		ToAccountID:   transfer1.ToAccountID,
		Amount:        transfer1.Amount,
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer2.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer2.Amount, args.Amount)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createTransferTest(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTransferTest(t)
	}
	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
