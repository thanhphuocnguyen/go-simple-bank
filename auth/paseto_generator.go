package auth

import (
	"log"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoGenerator struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoGenerator() (TokenGenerator, error) {
	generator := &PasetoGenerator{paseto.NewV4SymmetricKey(), []byte("my implicit nonce")}
	return generator, nil
}

func (g *PasetoGenerator) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	token := paseto.NewToken()
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	token.SetString("username", payload.Username)
	token.SetString("id", payload.ID.String())
	return token.V4Encrypt(g.symmetricKey, g.implicit), nil
}

func (g *PasetoGenerator) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.IssuedBy("simple-bank"))
	parsedToken, err := parser.ParseV4Local(g.symmetricKey, token, g.implicit)

	if err != nil {
		log.Println(err)
		if paseto.RuleError.Is(paseto.RuleError{}, err) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromParsedData(parsedToken)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func getPayloadFromParsedData(t *paseto.Token) (*Payload, error) {
	username, err := t.GetString("username")
	if err != nil {
		return nil, err
	}
	id, err := (t.GetString("id"))
	if err != nil {
		return nil, err
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        idUUID,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
