// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64              `json:"id"`
	Owner     string             `json:"owner"`
	Balance   int64              `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Entry struct {
	ID         int64 `json:"id"`
	IDAccounts int64 `json:"id_accounts"`
	// can be negative or positive
	Amount int64 `json:"amount"`
}

type Transfer struct {
	ID             int64 `json:"id"`
	FromIDAccounts int64 `json:"from_id_accounts"`
	ToIDAccounts   int64 `json:"to_id_accounts"`
	// must be positive
	Amount int64 `json:"amount"`
}
