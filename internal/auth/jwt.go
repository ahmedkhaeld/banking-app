package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretkey string
}

func NewJWTMaker() (*JWTMaker, error) {
	sk := os.Getenv("JWT_SECRET_KEY")
	if len(sk) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: %d must be at least %d char", len(sk), minSecretKeySize)
	}

	return &JWTMaker{
		secretkey: sk,
	}, nil
}

func (maker *JWTMaker) CreateToken(userId string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", payload, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(maker.secretkey))
	return tokenString, payload, err
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}

		return []byte(maker.secretkey), nil
	}
	clms := &Payload{}

	token, err := jwt.ParseWithClaims(tokenString, clms, keyFunc)
	if err != nil {
		ErrUnexpectedMethod = err
		return nil, ErrUnexpectedMethod
	}

	claims, ok := token.Claims.(*Payload)
	if !ok {
		ErrInvalidToken = err
		return nil, ErrInvalidToken
	}
	return claims, nil

}
