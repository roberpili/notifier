package session

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName = "session_token"
	TokenExpiration   = 6 * time.Hour
)

var (
	ErrNoCookie = errors.New("session cookie not found")
)

func SetAuthCookie(c *gin.Context, token string, config CookieConfig) {
	c.SetCookie(
		config.Name,
		token,
		int(TokenExpiration.Seconds()),
		"/",
		"", // empty for current domain
		config.Secure,
		config.HttpOnly,
	)
}

func GetAuthCookie(c *gin.Context, cookieName string) (string, error) {
	token, err := c.Cookie(SessionCookieName)
	if err != nil {
		return "", ErrNoCookie
	}

	return token, nil
}
