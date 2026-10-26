package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIDcontext, exists := c.Get("user_id")
	if !exists {
		slog.InfoContext(c.Request.Context(), "failed to get user_id from context")
		return 0, ErrUserIDNotFound
	}

	userID, ok := userIDcontext.(int)
	if !ok {
		slog.InfoContext(c.Request.Context(), "failed to cast user_id from context")
		return 0, ErrInvalidUserID
	}

	if userID <= 0 {
		slog.InfoContext(c.Request.Context(), "invalid user_id")
		return 0, ErrInvalidUserID
	}

	return userID, nil
}

func GetRoleFromContext(c *gin.Context) (entity.Role, error) {
	val, exists := c.Get("role")
	if !exists {
		return "", ErrMissingClaims
	}

	role, ok := val.(entity.Role)
	if !ok {
		return "", ErrMissingClaims
	}

	if role != entity.RoleAdmin && role != entity.RoleUser {
		return "", ErrInvalidRole
	}

	return role, nil
}
