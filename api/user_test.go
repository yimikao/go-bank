package api

import (
	db "gobank/db/sqlc"
	"gobank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPwd, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPwd,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}
