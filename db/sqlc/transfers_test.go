package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func createNewTransfer(t *testing.T) Transfer {
	account, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), 2)
	require.NoError(t, err)

	entryModel := CreateTransferParams{
		Amount:        utils.RandomMoney(),
		FromAccountID: account.ID,
		ToAccountID:   account2.ID,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), entryModel)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotEmpty(t, transfer.FromAccountID)
	require.NotEmpty(t, transfer.ToAccountID)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	createNewTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	entry1 := createNewTransfer(t)
	entry2, err := testQueries.GetTransfer(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.FromAccountID, entry2.FromAccountID)
	require.Equal(t, entry1.ToAccountID, entry2.ToAccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createNewTransfer(t)
	}

	arg := ListTransfersParams{
		FromAccountID: 1,
		ToAccountID:   1,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
