package vo

import (
	"errors"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type UpdateStatusVO struct {
	Status entity.Status
}

func NewUpdateStatusVO(status entity.Status) (*UpdateStatusVO, error) {
	if status != entity.StatusPending && status != entity.StatusCompleted {
		return nil, errors.New("status must be pending or completed")
	}
	return &UpdateStatusVO{Status: status}, nil
}
