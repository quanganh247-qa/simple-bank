package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {

	pwd := util.RandomString(6)

	hashedPwd1, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd1)

	err = CheckPassword(pwd, hashedPwd1)
	require.NoError(t, err)

	wrongPWd := util.RandomString(6)
	err = CheckPassword(wrongPWd, hashedPwd1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPwd2, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd2)
	require.NotEqual(t, hashedPwd1, hashedPwd2)
}
