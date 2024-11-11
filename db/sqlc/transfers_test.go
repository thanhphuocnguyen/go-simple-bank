package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func createNewTransfer(t *testing.T, account1ID, account2ID int64) Transfer {

	entryModel := CreateTransferParams{
		Amount:        utils.RandomMoney(),
		FromAccountID: account1ID,
		ToAccountID:   account2ID,
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
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createNewTransfer(t, account1.ID, account2.ID)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	entry1 := createNewTransfer(t, account1.ID, account2.ID)
	entry2, err := testQueries.GetTransfer(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.FromAccountID, entry2.FromAccountID)
	require.Equal(t, entry1.ToAccountID, entry2.ToAccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, 0)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createNewTransfer(t, account1.ID, account2.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
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
