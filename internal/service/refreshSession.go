package service

import (
	"context"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository/postgre"
)

const (
	maxUserSessions = 5
)

type RefreshSession struct {
	refreshSessionRepository postgre.RefreshSessionRepository
}

func NewRefreshSession(refreshSessionRepository postgre.RefreshSessionRepository) *RefreshSession {
	return &RefreshSession{
		refreshSessionRepository: refreshSessionRepository}
}

func (r *RefreshSession) PopSession(ctx context.Context, refreshToken string) (*model.RefreshSession, error) {
	session, err := r.refreshSessionRepository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	err = r.refreshSessionRepository.Delete(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *RefreshSession) SaveSession(ctx context.Context, session *model.RefreshSession) error {
	count, err := r.refreshSessionRepository.Count(ctx, session.UserID)

	if err != nil {
		return err
	}

	if count > maxUserSessions {
		errDelete := r.refreshSessionRepository.DeleteUserSessions(ctx, session.UserID)
		if errDelete != nil {
			return errDelete
		}
	}

	err = r.refreshSessionRepository.Create(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshSession) DeleteUserSessions(ctx context.Context, userId string) error {
	return r.refreshSessionRepository.DeleteUserSessions(ctx, userId)
}

func (r *RefreshSession) Delete(ctx context.Context, refreshToken string) error {
	return r.refreshSessionRepository.Delete(ctx, refreshToken)
}
