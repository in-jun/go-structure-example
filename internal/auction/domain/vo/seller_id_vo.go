package vo

import "github.com/in-jun/go-structure-example/internal/shared/validation"

type SellerIDVO struct{ ID string }

func NewSellerIDVO(id string) (*SellerIDVO, error) {
	parsed, err := validation.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return &SellerIDVO{ID: parsed}, nil
}
