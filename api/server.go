package api

import (
	"database/sql"
	"fmt"
	db "gobank/db/sqlc"
	"gobank/token"
	"gobank/util"
	"net/http"

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
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	r.POST("/users", server.createUser)
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts/", server.listAccount)
	r.POST("/transfers", server.createTransfer)
	server.router = r
	return
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	acc, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if acc.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", acc.ID, acc.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
