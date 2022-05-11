package db

import (
	"context"
	"testing"
	"time"

	"github.com/CM-IV/mef-api/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{

		UserName:       util.RandomOwner(),
		HashedPassword: "hashed secret", //Implement bcrypt hashing later!!!
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {

	createRandomUser(t)

}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}

func TestListUsers(t *testing.T) {

	for i := 0; i < 10; i++ {

		createRandomUser(t)

	}

	args := ListUsersParams{

		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {

		require.NotEmpty(t, user)

	}

}
