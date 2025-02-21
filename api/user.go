package api

import (
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type loginUserReq struct {
	Username string `json:"username" binding:"required,alphanum,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserReq struct {
	Username string `json:"username" binding:"required,alphanum,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type loginUserRes struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  userRes   `json:"user"`
}

type userRes struct {
	Username  string             `json:"username"`
	Fullname  string             `json:"fullname"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Sprintf("user %v not found", req.Username)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse("password is invalid"))
		return
	}

	accessToken, accessPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: [16]byte(refreshPayload.ID), Valid: true},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	res := loginUserRes{
		SessionID:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userRes{
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.Fullname,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			var message string

			switch pgErr.ConstraintName {
			case "users_pkey":
				message = fmt.Sprintf("username %v has been taken", req.Username)
			case "users_email_key":
				message = fmt.Sprintf("email %v has been taken", req.Email)
			}

			ctx.JSON(http.StatusForbidden, errorResponse(message))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	res := userRes{
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, res)
}
