package http

import (
	"testing"

	"github.com/in-jun/go-structure-example/internal/auth/application/command"
)

func TestToLoginResponse(t *testing.T) {
	result := &command.LoginResult{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	resp := toLoginResponse(result)

	if resp.AccessToken != "access-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "access-token")
	}
	if resp.RefreshToken != "refresh-token" {
		t.Errorf("RefreshToken = %q, want %q", resp.RefreshToken, "refresh-token")
	}
	if resp.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want 3600", resp.ExpiresIn)
	}
}

func TestToRefreshResponse(t *testing.T) {
	result := &command.RefreshResult{
		AccessToken:  "new-access",
		RefreshToken: "new-refresh",
		ExpiresIn:    1800,
	}

	resp := toRefreshResponse(result)

	if resp.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "new-access")
	}
	if resp.ExpiresIn != 1800 {
		t.Errorf("ExpiresIn = %d, want 1800", resp.ExpiresIn)
	}
}
