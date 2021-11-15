package employe

import (
	"context"
	"errors"
)

var (
	ErrEmployeNameNotSet = errors.New("employe name not set")
	ErrEmployeJobNotSet  = errors.New("employe job not set")
	ErrEmployeNotFound   = errors.New("employe not found")
	ErrCmdRepository     = errors.New("unable to command repository")
	ErrQueryRepository   = errors.New("unable to query repository")
)

// Service describes the Employe service.
type Service interface {
	Create(ctx context.Context, e EmployeDB) (string, error)
	GetByID(ctx context.Context, id string) (Employe, error)
	Status(ctx context.Context) (string, error)
}
