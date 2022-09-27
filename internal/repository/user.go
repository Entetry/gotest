package repository

import (
	"context"
	"entetry/gotest/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, company *model.User) (uuid.UUID, error)
	Update(ctx context.Context, company *model.User) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
