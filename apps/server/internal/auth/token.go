package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenManager struct {
	secret    []byte
	expiresIn time.Duration
	now       func() time.Time
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewTokenManager(secret string, expiresIn time.Duration) (*TokenManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if expiresIn <= 0 {
		return nil, errors.New("jwt expiration must be positive")
	}

	return &TokenManager{
		secret:    []byte(secret),
		expiresIn: expiresIn,
		now:       time.Now,
	}, nil
}

func (m *TokenManager) ExpiresInSeconds() int64 {
	return int64(m.expiresIn.Seconds())
}

func (m *TokenManager) Generate(userID uuid.UUID) (string, error) {
	now := m.now().UTC()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signed, nil
}

func (m *TokenManager) Parse(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithTimeFunc(m.now), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid subject", ErrInvalidToken)
	}

	return userID, nil
}
