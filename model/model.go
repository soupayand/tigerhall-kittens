package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int64
	Username string
	Password string
	Email    string
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email-id"`
	jwt.StandardClaims
}
