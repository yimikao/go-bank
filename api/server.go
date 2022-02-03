package api

import (
	"fmt"
	db "gobank/db/sqlc"
	"gobank/token"
	"gobank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(cfg util.Config, s db.Store) (server *Server, err error) {
	tm, err := token.NewJWTMaker(cfg.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server = &Server{
		config:     cfg,
		store:      s,
		tokenMaker: tm,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
func (s *Server) setupRouter() {
	r := gin.Default()

	r.POST("/users", s.createUser)
	r.POST("/users/login", s.login)

	authRoutes := r.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts/", s.listAccount)
	authRoutes.POST("/transfers", s.createTransfer)

	s.router = r
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
