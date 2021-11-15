package employehttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/mozgunovdm/example/internal/pkg/employe"
	"github.com/mozgunovdm/example/internal/pkg/transport"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

func NewService(
	svcEndpoints transport.Endpoints, logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
	}
	// HTTP Post - /employes
	r.Methods("POST").Path("/employes").Handler(kithttp.NewServer(
		svcEndpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /employes/{id}
	r.Methods("GET").Path("/employes/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /status
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
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
