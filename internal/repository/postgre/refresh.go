package postgre

import (
	"context"
	"crypto/sha1"
	"entetry/gotest/internal/auth/repository"
	"entetry/gotest/internal/config"
	"entetry/gotest/internal/model"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	ErrTokenInvalidOrExpired = errors.New("token is invalid or expired")
)

type Refresh struct {
	config            *config.Config
	refreshRepository *repository.Refresh
}

func NewRefresh(config *config.Config, refreshRepository *repository.Refresh) *Refresh {
	return &Refresh{
		config:            config,
		refreshRepository: refreshRepository,
	}
}

func (r *Refresh) GenerateAccessToken(userId int) (string, error) {
	claim := model.Claim{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(r.config.AccessTokenExpiration).Unix(),
		},
		UserID: userId,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).
		SignedString([]byte(r.config.AccessTokenKey))
}

func (r *Refresh) GenerateRefreshToken(userId int) (string, error) {
	claim := model.Claim{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(r.config.RefreshTokenExpiration).Unix(),
		},
		UserID: userId,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).
		SignedString([]byte(r.config.RefreshTokenKey))
}

func (r *Refresh) ValidateToken(ctx context.Context, token string) (int, error) {
	claims := new(model.Claim)

	parseToken, err := jwt.ParseWithClaims(token, claims, r.defaultKeyFunc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse: %v", err)
	}
	if !parseToken.Valid {
		return 0, ErrTokenInvalidOrExpired
	}

	hash, err := r.refreshRepository.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return 0, err
	}

	verifyHash := r.makeHash(token)

	if hash != verifyHash {
		return 0, ErrTokenInvalidOrExpired
	}

	return claims.UserID, nil
}

func (r *Refresh) defaultKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(r.config.RefreshTokenKey), nil
}

func (r *Refresh) makeHash(token string) string {
	h := sha1.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum([]byte(r.config.Salt)))
}
