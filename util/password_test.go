package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	match, err := VerifyPassword(password, hashedPassword1)
	require.NoError(t, err)
	require.True(t, match)

	wrongPassword := RandomString(6)
	match, err = VerifyPassword(wrongPassword, hashedPassword1)
	require.False(t, match)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
