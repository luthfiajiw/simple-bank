package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	Router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts", server.listAccounts)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.POST("/api/v1/transfer", server.createTransfer)

	server.Router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err string) gin.H {
	return gin.H{"error": err}
}
