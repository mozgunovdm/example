package employehttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mozgunovdm/example/employe"
	"github.com/mozgunovdm/example/employe/transport"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

// // NewService wires Go kit endpoints to the HTTP transport.
// func NewService(
// 	svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger log.Logger,
// ) http.Handler {
// 	// set-up router and initialize http endpoints
// 	var (
// 		r            = mux.NewRouter()
// 		errorLogger  = kithttp.ServerErrorLogger(logger)
// 		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
// 	)
// 	options = append(options, errorLogger, errorEncoder)
// 	//options := []kithttp.ServerOption{
// 	//	kithttp.ServerErrorLogger(logger),
// 	//	kithttp.ServerErrorEncoder(encodeError),
// 	//}
// 	// HTTP Post - /orders
// 	r.Methods("POST").Path("/employes").Handler(kithttp.NewServer(
// 		svcEndpoints.Create,
// 		decodeCreateRequest,
// 		encodeResponse,
// 		options...,
// 	))

// 	// HTTP Post - /orders/{id}
// 	r.Methods("GET").Path("/employes/{id}").Handler(kithttp.NewServer(
// 		svcEndpoints.GetByID,
// 		decodeGetByIDRequest,
// 		encodeResponse,
// 		options...,
// 	))

// 	return r
// }

func NewService(
	svcEndpoints transport.Endpoints, logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
	}
	// HTTP Post - /orders
	r.Methods("POST").Path("/employes").Handler(kithttp.NewServer(
		svcEndpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /orders/{id}
	r.Methods("GET").Path("/employes/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/status").Handler(kithttp.NewServer(
		svcEndpoints.Status,
		decodeStatusRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req transport.CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Employe); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return transport.GetByIDRequest{ID: id}, nil
}

func decodeStatusRequest(
	_ context.Context,
	r *http.Request,
) (request interface{}, err error) {
	return transport.StatusRequest{}, nil
}

// func encodeResponse(
// 	_ context.Context,
// 	w http.ResponseWriter,
// 	response interface{},
// ) error {
// 	return json.NewEncoder(w).Encode(response)
// }

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case employe.ErrEmployeNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
