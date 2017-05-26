package lorem

import (
	"context"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"encoding/json"
	"errors"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/kit/log"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func decodeLoremRequest(_ context.Context, r *http.Request) (interface{},error) {
	vars := mux.Vars(r)
	requestType,ok := vars["type"]

	if !ok {
		return nil,ErrBadRouting
	}

	vmin, ok := vars["min"]
	if !ok {
		return nil, ErrBadRouting
	}

	vmax, ok := vars["max"]
	if !ok {
		return nil, ErrBadRouting
	}

	min,_ :=strconv.Atoi(vmin)
	max,_ := strconv.Atoi(vmax)

	return LoremRequest {
		RequestType:requestType,
		Min:min,
		Max:max,
	},nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter,response interface{}) error {
	if e,ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx,e.error(),w)
		return nil
	}
	w.Header().Set("Context-type","application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	w.Header().Set("Content-type","application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func MakeHttpHandler(ctx context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption {
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	//POST /lorem/{type}/{min}/{max}
	r.Methods("POST").
		Path("/lorem/{type}/{min}/{max}").
		Handler(httptransport.NewServer(
		endpoint.LoremEndpoint,
		decodeLoremRequest,
		encodeResponse,
		options...,
	))

	return r
}