package vo

import "errors"

var (
	errCurrentPasswordRequired = errors.New("current password is required")
	errNewPasswordRequired     = errors.New("new password is required")
	errNewPasswordTooShort     = errors.New("new password must be at least 6 characters")
	errPasswordSame            = errors.New("new password must be different from current password")
)

type UpdatePasswordVO struct {
	CurrentPassword string
	NewPassword     string
}

func NewUpdatePasswordVO(currentPassword, newPassword string) (*UpdatePasswordVO, error) {
	if currentPassword == "" {
		return nil, errCurrentPasswordRequired
	}
	if newPassword == "" {
		return nil, errNewPasswordRequired
	}
	if len(newPassword) < 6 {
		return nil, errNewPasswordTooShort
	}
	if currentPassword == newPassword {
		return nil, errPasswordSame
	}
	return &UpdatePasswordVO{CurrentPassword: currentPassword, NewPassword: newPassword}, nil
}
