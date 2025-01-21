package db_test

import (
	"context"
	db "simplebank/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := db.NewStore(testPool)

	account1 := createTestAcccount(t)
	account2 := createTestAcccount(t)

	// Run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan db.TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromIDAccounts: account1.ID,
				ToIDAccounts:   account2.ID,
				Amount:         amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results & errs
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromIDAccounts, account1.ID)
		require.Equal(t, transfer.ToIDAccounts, account2.ID)
		require.Equal(t, transfer.Amount, amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check Entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.IDAccounts, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.IDAccounts, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

	}

}
