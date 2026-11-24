package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

type Claims struct {
	Email  string      `json:"email"`
	UserID int         `json:"user_id"`
	Role   entity.Role `json:"role"`
	jwt.RegisteredClaims
}

type GoogleUser struct {
	Sub   string
	Email string
	Name  string
}
