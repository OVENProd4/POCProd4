package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"Book-repo/internal/util"
	"Book-repo/pkg/database/endpoints"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	//	m := http.NewServeMux()
	r := mux.NewRouter()

	r.Methods("GET").Path("/healthz").Handler(httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/update/{ticketid,pagenumber}").Handler(httptransport.NewServer(
		ep.UpdateEndpoint,
		decodeHTTPUpdateRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/add/{doc}").Handler(httptransport.NewServer(
		ep.AddEndpoint,
		decodeHTTPAddRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/get/{key,value}").Handler(httptransport.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/remove").Handler(httptransport.NewServer(
		ep.RemoveEndpoint,
		decodeHTTPRemoveRequest,
		encodeResponse,
	))

	return r
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRequest
	if r.ContentLength == 0 {
		logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	/*	vars := m.Vars(r)
		team, ok := vars["team"]
		if !ok {
			return nil, ErrBadRouting
		} */
	var req endpoints.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPRemoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
