package db_test

import (
	"context"
	"fmt"
	db "simplebank/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := db.NewStore(testPool)

	account1 := createTestAcccount(t)
	account2 := createTestAcccount(t)
	fmt.Println(">> Balance Before:", account1.Balance, account2.Balance)

	// Run n concurrent transfer transactions
	n := 2
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

		// Check Accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// Check balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2*amount, 3*amount, ..., n*amount
	}

	// Check the final updated balances
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> Balance After:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
