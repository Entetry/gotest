package handlers

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string
	RefreshToken string
}
