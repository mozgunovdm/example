package implementation

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	"github.com/rickb777/date"

	employeService "github.com/mozgunovdm/example/employe"
)

// service implements the Employe Service
type service struct {
	repository employeService.Repository
	logger     log.Logger
}

// NewService creates and returns a new service instance
func NewService(rep employeService.Repository, logger log.Logger) employeService.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Create employe
func (s *service) Create(ctx context.Context, employe employeService.Employe) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	level.Info(logger).Log("Create", 1)
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	employe.ID = id
	employe.Name = "Employe"
	employe.Job = "Robot"
	employe.EmployedAt = date.Today().FormatISO(4)

	return id, nil
}

// GetByID returns an employe given by id
func (s *service) GetByID(ctx context.Context, id string) (employeService.Employe, error) {
	logger := log.With(s.logger, "method", "GetByID")
	employe, err := s.repository.GetEmployeByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return employe, employeService.ErrEmployeNotFound
		}
		return employe, employeService.ErrQueryRepository
	}
	return employe, nil
}

//Status return info that service ok
func (_ *service) Status(ctx context.Context) (string, error) {
	return "ok", nil
}
