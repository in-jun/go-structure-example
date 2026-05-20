package http

import "github.com/in-jun/go-structure-example/internal/auth/application/command"

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func toAuthResponse(r *command.LoginResult) *AuthResponse {
	return &AuthResponse{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		ExpiresIn:    r.ExpiresIn,
	}
}

func toRefreshResponse(r *command.RefreshResult) *AuthResponse {
	return &AuthResponse{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		ExpiresIn:    r.ExpiresIn,
	}
}
