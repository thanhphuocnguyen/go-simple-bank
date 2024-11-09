package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func createNewEntry(t *testing.T) Entry {
	account, err := testQueries.GetAccount(context.Background(), 1)

	require.NoError(t, err)

	entryModel := CreateEntryParams{
		Amount:    utils.RandomMoney(),
		AccountID: account.ID,
	}
	entry, err := testQueries.CreateEntry(context.Background(), entryModel)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return entry
}
func TestCreateEntry(t *testing.T) {
	createNewEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createNewEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createNewEntry(t)
	}

	arg := ListEntriesParams{
		AccountID: 1,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
