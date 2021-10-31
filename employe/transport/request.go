package transport

import (
	"github.com/mozgunovdm/example/employe"
)

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	Employe employe.Employe
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Employe employe.Employe `json:"employe"`
	Err     error           `json:"error,omitempty"`
}

type StatusRequest struct{}

type StatusResponse struct {
	Status string `json:"status"`
}
