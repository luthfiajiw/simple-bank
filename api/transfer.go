package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type transferRequest struct {
	FromIDAccounts int64  `json:"from_id_accounts" binding:"required,min=1"`
	ToIDAccounts   int64  `json:"to_id_accounts" binding:"required,min=1"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	Currency       string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	fromAccount, isValid := server.isCurrencyValid(ctx, req.FromIDAccounts, req.Currency)
	if !isValid {
		return
	}

	if authPayload.Username != fromAccount.Owner {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err.Error()))
		return
	}

	_, isValid = server.isCurrencyValid(ctx, req.ToIDAccounts, req.Currency)
	if !isValid {
		return
	}

	arg := db.TransferTxParams{
		FromIDAccounts: req.FromIDAccounts,
		ToIDAccounts:   req.ToIDAccounts,
		Amount:         req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (server *Server) isCurrencyValid(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err.Error()))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d has different currency", accountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return account, false
	}

	return account, true
}
