package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
