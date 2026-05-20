package vo

import "errors"

type UpdateProfileVO struct {
	Name string
}

func NewUpdateProfileVO(name string) (*UpdateProfileVO, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(name) > 100 {
		return nil, errors.New("name must be 100 characters or less")
	}
	return &UpdateProfileVO{Name: name}, nil
}
