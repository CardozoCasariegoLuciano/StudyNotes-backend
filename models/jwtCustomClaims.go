package models

import "github.com/dgrijalva/jwt-go"

type JwtCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}
