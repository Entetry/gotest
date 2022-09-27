package service

import (
	"context"
	"entetry/gotest/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	userService User
}

func NewAuthService(userService User) *Auth {
	return &Auth{
		userService: userService}
}

func (a Auth) AttemptLogin(context context.Context, username, password string) (bool, error) {
	user, err := a.userService.GetByUsername(context, username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a Auth) Register(context context.Context, user *model.User) (interface{}, interface{}) {

}

func comparePassword(password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
