package auth

import "time"

type TokenGenerator interface {
	GenerateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
