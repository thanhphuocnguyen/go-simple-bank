package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thanhphuocnguyen/go-simple-bank/auth"
	db "github.com/thanhphuocnguyen/go-simple-bank/db/sqlc"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

type Server struct {
	config         utils.Config
	store          db.Store
	router         *gin.Engine
	tokenGenerator auth.TokenGenerator
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := auth.NewPasetoGenerator()
	if err != nil {
		return nil, err
	}

	server := &Server{store: store, tokenGenerator: tokenGenerator, config: config}
	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// add routes for user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// add routes for accounts
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)

	// add routes for transfers
	router.POST("/transfers", server.createTransfer)
	server.router = router
}
