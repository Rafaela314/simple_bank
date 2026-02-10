package db

import (
	"context"
	"testing"
	"time"

	"simple_bank/util"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

// createRandomUser creates a random user for testing
func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: util.RandomString(6),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestGetUser tests the GetUser function
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

// TestGetUserNotFound tests the GetUser function when the user is not found.
func TestGetUserNotFound(t *testing.T) {
	_, err := testQueries.GetUser(context.Background(), "nonexistent_username")
	require.Error(t, err)
	require.Equal(t, pgx.ErrNoRows.Error(), err.Error())
}
