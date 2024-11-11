package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	randomPassword := RandomString(8)

	hashedPassword, err := HashPassword(randomPassword)
	require.NoError(t, err)
	require.NotEmptyf(t, hashedPassword, "hashed password should not be empty")

	err = ComparePassword(randomPassword, hashedPassword)
	require.NoError(t, err)
	wrongPassword := RandomString(8)
	err = ComparePassword(wrongPassword, hashedPassword)
	require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPassword2, err := HashPassword(randomPassword)
	require.NoError(t, err)
	require.NotEmptyf(t, hashedPassword2, "hashed password should not be empty")
	require.NotEqual(t, hashedPassword, hashedPassword2)
}
