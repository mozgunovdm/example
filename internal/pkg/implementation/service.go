package implementation

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/rickb777/date"

	"github.com/mozgunovdm/example/internal/pkg/employe"
)

// service implements the Employe Service
type empService struct {
	repository employe.Repository
	logger     log.Logger
}

// NewService creates and returns a new service instance
func NewService(rep employe.Repository, l log.Logger) employe.Service {
	return &empService{
		repository: rep,
		logger:     l,
	}
}

// Create employe
func (s *empService) Create(ctx context.Context, emp employe.EmployeDB) (string, error) {
	l := log.With(s.logger, "method", "Create")
	level.Info(l).Log("Create", 1)

	id := ""
	if emp.Name == "" {
		return id, employe.ErrEmployeNameNotSet
	}
	if emp.Job == "" {
		return id, employe.ErrEmployeJobNotSet
	}
	if emp.EmployedAt == "" {
		emp.EmployedAt = date.Today().String()
	}

	id, err := s.repository.CreateEmploye(ctx, emp)
	if err != nil {
		level.Error(l).Log("err", err)
		if err == sql.ErrNoRows {
			return id, employe.ErrEmployeNotFound
		}
		return id, employe.ErrQueryRepository
	}
	return id, nil
}

// GetByID returns an employe given by id
func (s *empService) GetByID(ctx context.Context, id string) (employe.Employe, error) {
	l := log.With(s.logger, "method", "GetByID")
	emp, err := s.repository.GetEmployeByID(ctx, id)
	if err != nil {
		level.Error(l).Log("err", err)
		if err == sql.ErrNoRows {
			return emp, employe.ErrEmployeNotFound
		}
		return emp, employe.ErrQueryRepository
	}
	return emp, nil
}

//Status return info that service ok
func (_ *empService) Status(ctx context.Context) (string, error) {
	return "ok", nil
}
