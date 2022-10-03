package postgre

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/model"
)

type RefreshSessionRepository interface {
	Create(ctx context.Context, session *model.RefreshSession) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshSession, error)
	Count(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, refreshToken string) error
	DeleteUserSessions(ctx context.Context, userID string) error
}

type RefreshSession struct {
	db *pgxpool.Pool
}

func NewRefresh(db *pgxpool.Pool) *RefreshSession {
	return &RefreshSession{db: db}
}

func (r *RefreshSession) GetHashByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	var hash string

	err := r.db.QueryRow(ctx, "SELECT hash FROM refresh_sessions WHERE user_id = $1", userID).Scan(&hash)
	if err != nil {
		log.Errorf("Get hash by user id failed : %v\n", err)
		return "", err
	}

	return hash, nil
}

func (r *RefreshSession) Create(ctx context.Context, session *model.RefreshSession) error {
	_, err := r.db.Exec(ctx, `INSERT INTO refresh_sessions (refresh_token, user_id, ua, fingerprint, ip, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, session.RefreshToken, session.UserID, session.UserAgent, session.Fingerprint,
		session.IP, session.ExpiresAt)
	if err != nil {
		log.Errorf("Cannot create an refresh session: %v\n", err)
		return fmt.Errorf("cannot create an refresh session: %v", err)
	}
	return nil
}

func (r *RefreshSession) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshSession, error) {
	var session model.RefreshSession
	err := r.db.QueryRow(ctx, `SELECT refresh_token, user_id, ua, fingerprint, ip, expires_at FROM refresh_sessions 
		WHERE refresh_token = $1`, refreshToken).Scan(&session.RefreshToken, &session.UserID, &session.UserAgent, &session.Fingerprint,
		&session.IP, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("cannot get refreshSession: %v", err)
	}
	return &session, nil
}

func (r *RefreshSession) Count(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT count(1) FROM refresh_sessions WHERE user_id = $1`, userID).Scan(&count)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("error in Count Refresh Session: %v", err)
	}
	return count, nil
}

func (r *RefreshSession) Delete(ctx context.Context, refreshToken string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_sessions WHERE refresh_token = $1`, refreshToken)
	if err != nil {
		return fmt.Errorf("can't DeleteSession in users: %v", err)
	}
	return nil
}

func (r *RefreshSession) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_sessions WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("can't delete user sessions: %v", err)
	}
	return nil
}
