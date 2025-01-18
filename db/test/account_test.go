package db_test

import (
	"context"
	db "simplebank/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestAcccount(t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		Owner: "luthfi",
		Balance: 0,
		Currency: "IDR",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T)  {
	createTestAcccount(t)
}

func TestGetAccount(t *testing.T)  {
	createdAccount := createTestAcccount(t)

	account, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdAccount.ID, account.ID)
	require.Equal(t, createdAccount.Owner, account.Owner)
	require.Equal(t, createdAccount.Balance, account.Balance)
	require.Equal(t, createdAccount.Currency, account.Currency)
}

