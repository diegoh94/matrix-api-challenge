package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (service *TokenService) CreateAccessToken() (string, int64, error) {
	expiresAt := time.Now().Add(service.ttl)

	claims := jwt.MapClaims{
		"sub": "matrix-api-client",
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(service.secret)
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return signedToken, int64(service.ttl.Seconds()), nil
}

func (service *TokenService) ValidateAccessToken(tokenString string) error {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return service.secret, nil
	})
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	if !parsedToken.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
