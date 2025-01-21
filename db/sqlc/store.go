package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	*pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(pool),
		Pool:    pool,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromIDAccounts int64 `json:"from_id_accounts"`
	ToIDAccounts   int64 `json:"to_id_accounts"`
	Amount         int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another.
// It creates a transfer record, add account entries and update account balance
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			IDAccounts: arg.FromIDAccounts,
			Amount:     -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			IDAccounts: arg.ToIDAccounts,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		// UPDATE ACCOUNT BALANCE
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     arg.FromIDAccounts,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     arg.ToIDAccounts,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
