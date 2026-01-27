package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl}
}

func (j *JWT) IssueAccessToken(userID string, role string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(j.ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

type AuthClaims struct {
	UserID string
	Role   string
}

func (j *JWT) Parse(tokenStr string) (AuthClaims, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return AuthClaims{}, errors.New("empty token")
	}

	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil || !parsed.Valid {
		return AuthClaims{}, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return AuthClaims{}, errors.New("invalid claims")
	}

	sub, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	if sub == "" {
		return AuthClaims{}, errors.New("missing sub")
	}

	return AuthClaims{UserID: sub, Role: role}, nil
}
