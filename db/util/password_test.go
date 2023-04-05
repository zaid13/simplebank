package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {

	password := RandomString(6)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return
	}
	err = CheckPassword(password, hashedPassword)

	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	wrongPassword = RandomString(2)
	err = CheckPassword(password, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrHashTooShort.Error())

}
