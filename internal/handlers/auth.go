package handlers

import (
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Auth struct {
	authService *service.Auth
}

func NewAuth(authService *service.Auth) *Auth {
	return &Auth{authService: authService}
}

func (a *Auth) SignIn(ctx echo.Context) error {
	request := new(signInRequest)
	if err := ctx.Bind(request); err != nil {
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

func (a *Auth) SignUp(ctx echo.Context) error {
	request := new(signUpRequest)
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := a.authService.SignUp(ctx.Request().Context(), request.Username, request.Password, request.Email)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusCreated, "Registration completed successfully")
}

func (a *Auth) Refresh(ctx echo.Context) error {
	request := new(refreshTokenRequest)
	if err := ctx.Bind(request); err != nil {
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

func (a *Auth) Logout(ctx echo.Context) error {
	request := new(logoutRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
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
