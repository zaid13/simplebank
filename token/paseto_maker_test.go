package token

import (
	//"github.com/dgrijalva/Paseto-go"
	//"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
	"github.com/zaid13/simplebank/db/util"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := CreatePasetoMakerInstance(util.RandomString(32))
	require.NoError(t, err)

	userNmae := util.RandomOwner()
	durartion := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(durartion)

	token, err := maker.CreateToken(userNmae, durartion)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userNmae, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {

	maker, err := CreatePasetoMakerInstance(util.RandomString(40))
	require.NoError(t, err)

	userNmae := util.RandomOwner()
	durartion := time.Minute

	token, err := maker.CreateToken(userNmae, -durartion)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())

	require.Nil(t, payload)

}
