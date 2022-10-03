package service

import (
	"context"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository/postgre"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepository postgre.UserRepository
}

func NewUserService(userRepository postgre.UserRepository) *User {
	return &User{
		userRepository: userRepository}
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.userRepository.GetByUsername(ctx, username)
}

func (u *User) Create(ctx context.Context, username, password, email string) (uuid.UUID, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return uuid.Nil, err
	}

	return u.userRepository.Create(ctx, username, string(pwdHash), email)
}
