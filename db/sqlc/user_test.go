package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zaid13/simplebank/db/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hash, err := util.HashPassword(password)
	arg := CreateUsersParams{
		Username:     util.RandomOwner(),
		HashPassword: hash,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}

	var user User
	user, err = testQueries.CreateUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	errr := util.CheckPassword(password, arg.HashPassword)
	require.NoError(t, errr)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	//crete account
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUsers(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashPassword, user2.HashPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)

}
