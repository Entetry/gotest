package service

import (
	"context"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository"
	"github.com/google/uuid"
)

type User struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *User {
	return &User{
		userRepository: userRepository}
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.userRepository.GetByUsername(ctx, username)
}

func (u *User) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	return u.userRepository.Create(ctx, user)
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	return u.userRepository.Update(ctx, user)
}

func (u *User) Delete(ctx context.Context, uuid uuid.UUID) error {
	return u.userRepository.Delete(ctx, uuid)
}
