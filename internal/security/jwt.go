package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService issues and verifies JWTs for WebSocket authentication.
type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), ttl: ttl}
}

type Claims struct {
	Subject string `json:"sub"`
	jwt.RegisteredClaims
}

func (s *TokenService) Issue(subject string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": subject,
		"iat": now.Unix(),
		"exp": now.Add(s.ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *TokenService) Verify(token string) (string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	sub, _ := claims.GetSubject()
	if sub == "" {
		return "", errors.New("missing subject")
	}
	return sub, nil
}