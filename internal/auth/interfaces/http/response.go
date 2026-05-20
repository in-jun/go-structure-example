package http

import "github.com/in-jun/go-structure-example/internal/auth/application/command"

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func toLoginResponse(r *command.LoginResult) *Response {
	return &Response{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		ExpiresIn:    r.ExpiresIn,
	}
}

func toRefreshResponse(r *command.RefreshResult) *Response {
	return &Response{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		ExpiresIn:    r.ExpiresIn,
	}
}
