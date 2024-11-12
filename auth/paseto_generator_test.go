package auth

import (
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func TestPasetoValidToken(t *testing.T) {
	generator, err := NewPasetoGenerator()
	require.NoError(t, err)

	randomUserName := utils.RandomOwner()
	duration := time.Minute
	token, err := generator.GenerateToken(randomUserName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, randomUserName, payload.Username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)
}

func TestPasetoExpiredToken(t *testing.T) {
	generator, err := NewPasetoGenerator()
	require.NoError(t, err)

	token, err := generator.GenerateToken(utils.RandomOwner(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestPasetoInvalidToken(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	pasetoToken := paseto.NewToken()
	pasetoToken.SetIssuedAt(payload.IssuedAt)
	pasetoToken.SetExpiration(payload.ExpiredAt)
	pasetoToken.SetString("username", payload.Username)
	pasetoToken.SetString("id", payload.ID.String())
	token := pasetoToken.V4Encrypt(paseto.NewV4SymmetricKey(), []byte("my implicit nonce"))
	require.NotEmpty(t, token)
	maker, err := NewPasetoGenerator()
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
