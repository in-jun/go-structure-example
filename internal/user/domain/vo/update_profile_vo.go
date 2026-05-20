package vo

import "errors"

var (
	errNameRequired = errors.New("name is required")
	errNameTooLong  = errors.New("name must be 100 characters or less")
)

type UpdateProfileVO struct {
	Name string
}

func NewUpdateProfileVO(name string) (*UpdateProfileVO, error) {
	if name == "" {
		return nil, errNameRequired
	}
	if len(name) > 100 {
		return nil, errNameTooLong
	}
	return &UpdateProfileVO{Name: name}, nil
}
