package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	pwd := RandomString(6)

	hashedPwd, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd)

	err = CheckPassword(pwd, hashedPwd)
	require.NoError(t, err)

	// wrongPassword := RandomString(6)
	// err = CheckPassword(wrongPassword, hashedPassword)
	// require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
