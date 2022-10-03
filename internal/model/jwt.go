package model

import "github.com/golang-jwt/jwt"

type RefreshSession struct {
	RefreshToken string
	UserID       string
	ExpiresAt    int64
	TokenParam
}

type TokenParam struct {
	UserAgent   string
	Fingerprint string
	IP          string
}

type Claim struct {
	UserID string
	jwt.StandardClaims
}
