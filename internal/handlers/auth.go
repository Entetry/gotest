package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/model"
	"entetry/gotest/internal/service"
)

// Auth handler struct
type Auth struct {
	authService *service.Auth
}

// NewAuth creates new auth handler
func NewAuth(authService *service.Auth) *Auth {
	return &Auth{authService: authService}
}

// SignIn sign in into account
func (a *Auth) SignIn(ctx echo.Context) error {
	request := new(signInRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokenParam, err := parseTokenParam(ctx.Request().Header)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	refreshToken, accessToken, err := a.authService.SignIn(ctx.Request().Context(), request.Username, request.Password, tokenParam)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, &tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// SignUp sign up into account
func (a *Auth) SignUp(ctx echo.Context) error {
	request := new(signUpRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = a.authService.SignUp(ctx.Request().Context(), request.Username, request.Password, request.Email)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusCreated, "Registration completed successfully")
}

// Refresh update refresh token
func (a *Auth) Refresh(ctx echo.Context) error {
	request := new(refreshTokenRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	tokenParam, err := parseTokenParam(ctx.Request().Header)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	refreshToken, accessToken, err := a.authService.RefreshTokens(ctx.Request().Context(), request.RefreshToken, tokenParam)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, &tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout log out from session
func (a *Auth) Logout(ctx echo.Context) error {
	request := new(logoutRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = a.authService.Logout(ctx.Request().Context(), request.RefreshToken)

	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func parseTokenParam(header http.Header) (*model.TokenParam, error) {
	ua := header.Get("User-Agent")
	fingerprint := header.Get("Fingerprint")
	ip := header.Get("IP")

	if ua == "" || fingerprint == "" || ip == "" {
		return nil, errors.New("parameters of header missing")
	}

	return &model.TokenParam{
		UserAgent:   ua,
		Fingerprint: fingerprint,
		IP:          ip,
	}, nil
}
