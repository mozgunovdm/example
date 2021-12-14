package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mozgunovdm/example/internal/pkg/employe"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	Status  endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s employe.Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
		Status:  makeStatusEndpoints(s),
	}
}

func makeCreateEndpoint(s employe.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest) // type assertion
		id, err := s.Create(ctx, req.Employe)
		return CreateResponse{ID: id, Err: err}, err
	}
}

func makeGetByIDEndpoint(s employe.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		employeRes, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Employe: employeRes, Err: err}, err
	}
}

func makeStatusEndpoints(s employe.Service) endpoint.Endpoint {
	return func(
		ctx context.Context,
		request interface{},
	) (interface{}, error) {
		res, err := s.Status(ctx)
		return StatusResponse{res}, err
	}
}
