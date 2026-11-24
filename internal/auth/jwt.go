package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

// authentication and authorization token
type TokenGenerator interface {
	Generate(userID int, email string, role entity.Role) (string, error)
	Validate(tokenString string) (*Claims, error)
}

type JWTGenerator struct {
	jwtKey               []byte
	tokenExpirationHours int
}

func NewJWTGenerator(jwtKey []byte, expirationHours int) *JWTGenerator {
	return &JWTGenerator{
		jwtKey:               jwtKey,
		tokenExpirationHours: expirationHours,
	}
}

func (j *JWTGenerator) Generate(userID int, email string, role entity.Role) (string, error) {
	expiration := time.Now().Add(time.Duration(j.tokenExpirationHours) * time.Hour)

	claims := &Claims{
		Email:  email,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Audience:  []string{"notifier"},
			Issuer:    "notifier-service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtKey)
}

func (j *JWTGenerator) Validate(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return j.jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("JWTGenerator.Validate failed to parse claims: %w", err)
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
