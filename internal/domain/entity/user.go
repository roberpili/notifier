package entity

import (
	"time"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	GoogleID     string    `json:"-"`
	Name         string    `json:"name"`
	IsActive     bool      `json:"isActive"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}
