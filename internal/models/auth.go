package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
