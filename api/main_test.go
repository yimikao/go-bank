package api

import (
	db "gobank/db/sqlc"
	"gobank/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, s db.Store) *Server {
	c := util.Config{
		TokenSecretKey:      util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	svr, err := NewServer(c, s)
	require.NoError(t, err)

	return svr
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
