package transport

import (
	"github.com/mozgunovdm/example/internal/pkg/employe"
)

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	Employe employe.EmployeDB
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Employe employe.Employe `json:"employe"`
	Err     error           `json:"error"`
}

type StatusRequest struct{}

type StatusResponse struct {
	Status string `json:"status"`
}
