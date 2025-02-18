package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "users":
				ctx.JSON(http.StatusForbidden, errorResponse(fmt.Sprintf("this owner already has an account with %v currency", req.Currency)))
				return
			case "accounts_owner_fkey":
				ctx.JSON(http.StatusForbidden, errorResponse(fmt.Sprintf("owner with username %v was not found", arg.Owner)))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type listAccountReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err.Error()))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	if account.Owner != authPayload.Username {
		err = errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
