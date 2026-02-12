package api

import (
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router.POST("/accounts", server.createAccountHandler)
	server.router.GET("/accounts/:id", server.getAccountHandler)
	server.router.GET("/accounts", server.listAccountsHandler)
	server.router.POST("/transfers", server.createTransferHandler)

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
