package session

import (
	"time"
)

type CookieConfig struct {
	Name            string
	TokenExpiration time.Duration
	Secure          bool // true in https
	HttpOnly        bool
	SameSite        int
}
