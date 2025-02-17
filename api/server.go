package api

import (
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	config     utils.Config
	TokenMaker token.Maker
	Router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SymetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		config:     config,
		TokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/api/v1/users", server.createUser)
	router.POST("/api/v1/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(AuthMiddleware(tokenMaker))

	authRoutes.POST("/api/v1/accounts", server.createAccount)
	authRoutes.GET("/api/v1/accounts", server.listAccounts)
	authRoutes.GET("/api/v1/accounts/:id", server.getAccount)

	authRoutes.POST("/api/v1/transfer", server.createTransfer)

	server.Router = router

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err string) gin.H {
	return gin.H{"error": err}
}
