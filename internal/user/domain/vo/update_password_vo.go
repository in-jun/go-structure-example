package vo

import "errors"

type UpdatePasswordVO struct {
	CurrentPassword string
	NewPassword     string
}

func NewUpdatePasswordVO(currentPassword, newPassword string) (*UpdatePasswordVO, error) {
	if currentPassword == "" {
		return nil, errors.New("current password is required")
	}
	if newPassword == "" {
		return nil, errors.New("new password is required")
	}
	if len(newPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters")
	}
	if currentPassword == newPassword {
		return nil, errors.New("new password must be different from current password")
	}
	return &UpdatePasswordVO{CurrentPassword: currentPassword, NewPassword: newPassword}, nil
}
