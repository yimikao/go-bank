package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	pwd := RandomString(6)

	hashedPwd1, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd1)

	err = CheckPassword(pwd, hashedPwd1)
	require.NoError(t, err)

	// wrongPassword := RandomString(6)
	// err = CheckPassword(wrongPassword, hashedPassword)
	// require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPwd2, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd2)
	require.NotEqual(t, hashedPwd1, hashedPwd2)
}
