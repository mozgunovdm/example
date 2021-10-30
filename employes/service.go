package employe

import (
	"context"
	"errors"
)

var (
	ErrEmployeNotFound = errors.New("employe not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

// Service describes the Employe service.
type Service interface {
	Create(ctx context.Context, order Employe) (string, error)
	GetByID(ctx context.Context, id string) (Employe, error)
}
