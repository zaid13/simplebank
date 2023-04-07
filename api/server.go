package api

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/zaid13/simplebank/db/sqlc"
	"github.com/zaid13/simplebank/db/util"
	"github.com/zaid13/simplebank/token"
)
import "github.com/gin-gonic/gin"

type Server struct {
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {

		return nil, fmt.Errorf("cannot create the token maker %w", err)

	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", ValidCurrency)
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)


	 authRoutes:=router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)


	server.router = router
}
