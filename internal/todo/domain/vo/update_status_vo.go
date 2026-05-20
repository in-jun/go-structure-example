package vo

import (
	"errors"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

var errInvalidStatus = errors.New("status must be pending or completed")

type UpdateStatusVO struct {
	Status entity.Status
}

func NewUpdateStatusVO(status entity.Status) (*UpdateStatusVO, error) {
	if status != entity.StatusPending && status != entity.StatusCompleted {
		return nil, errInvalidStatus
	}
	return &UpdateStatusVO{Status: status}, nil
}
