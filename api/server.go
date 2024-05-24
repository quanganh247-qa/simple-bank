package api

import (
	"github.com/gin-gonic/gin"
	db "tutorial.sqlc.dev/app/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Starts run the HTTP server on a specific address
func (server *Server) Starts(address string) error {
	return server.router.Run(address)
}

func errorMessage(err error) gin.H {
	return gin.H{"error": err.Error()}
}
