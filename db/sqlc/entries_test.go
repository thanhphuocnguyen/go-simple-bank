package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func createNewEntry(t *testing.T, accountID int64) Entry {
	account, err := testQueries.GetAccount(context.Background(), accountID)
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
	account := createRandomAccount(t)
	createNewEntry(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)

	entry1 := createNewEntry(t, account.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, 0)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createNewEntry(t, account.ID)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
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
