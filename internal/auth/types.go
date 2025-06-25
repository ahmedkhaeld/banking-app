package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrUnexpectedMethod = errors.New("unexpected signing method")
)

type TokenMaker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	jwt.MapClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{MapClaims: jwt.MapClaims{
		"iss": tokenID,
		"sub": username,
		"aud": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
		"nbf": time.Now().Unix(),
		"jti": uuid.New().String(),
	},
	}, nil
}

func (payload *Payload) Valid() error {
	expiration, err := payload.GetExpirationTime()
	if err != nil {
		return err
	}
	if time.Now().After(expiration.Time) {
		return ErrExpiredToken
	}
	return nil
}
