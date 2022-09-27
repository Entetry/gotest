package handlers

import (
	"entetry/gotest/internal/handlers/request"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository/postgre"
	"entetry/gotest/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Auth struct {
	authService    *service.Auth
	refreshService *postgre.Refresh
}

func NewAuth(authService *service.Auth) *Auth {
	return &Auth{authService: authService}
}

func (a *Auth) SignIn(ctx echo.Context) error {
	request := new(request.RegisterRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ok, err := a.authService.AttemptLogin(ctx.Request().Context(), request.Username, request.Password)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	token := a.refreshService.GenerateAccessToken()
}

func (a *Auth) SignUp(ctx echo.Context) error {
	request := new(request.RegisterRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	var user *model.User
	user.Username = request.Username
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	user.PasswordHash = string(hash)
	id, err := a.authService.Register(ctx.Request().Context(), user)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, id)
}

func (a *Auth) Refresh(ctx echo.Context) error {

}

func (a *Auth) Logout(ctx echo.Context) error {

}
