package validation

import "github.com/google/uuid"

func ParseUUID(s string) (string, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return s, nil
}
