package transport

import (
	"Book-repo/pkg/book/endpoints"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"Book-repo/internal/util"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/healthz").Handler(httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/bookmark/{ticketid}/{pagenumber}").Handler(httptransport.NewServer(
		ep.BookmarkEndpoint,
		decodeHTTPBookmarkRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/addDocument").Handler(httptransport.NewServer(
		ep.AddDocumentEndpoint,
		decodeHTTPAddDocumentRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/get").Handler(httptransport.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/status/{ticketid}").Handler(httptransport.NewServer(
		ep.StatusEndpoint,
		decodeHTTPStatusRequest,
		encodeResponse,
	))
	/*m := http.NewServeMux()

	m.Handle("/healthz", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	m.Handle("/status", httptransport.NewServer(
		ep.StatusEndpoint,
		decodeHTTPStatusRequest,
		encodeResponse,
	))
	m.Handle("/addDocument", httptransport.NewServer(
		ep.AddDocumentEndpoint,
		decodeHTTPAddDocumentRequest,
		encodeResponse,
	))
	m.Handle("/get", httptransport.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	m.Handle("/bookmark", httptransport.NewServer(
		ep.BookmarkEndpoint,
		decodeHTTPBookmarkRequest,
		encodeResponse,
	)) */

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

func decodeHTTPStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ticketid, ok := vars["ticketid"]
	if !ok {
		return nil, ErrBadRouting
	}
	/*var req endpoints.StatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	} */
	return endpoints.StatusRequest{TicketID: ticketid}, nil
}

func decodeHTTPBookmarkRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	ticketid, ok := vars["ticketid"]
	if !ok {
		return nil, ErrBadRouting
	}

	pagenumber, ok := vars["pagenumber"]
	if !ok {
		return nil, ErrBadRouting
	}
	/*var req endpoints.BookmarkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	} */
	i, _ := strconv.ParseInt(pagenumber, 0, 64)

	return endpoints.BookmarkRequest{TicketID: ticketid, PageNumber: int(i)}, nil
}

func decodeHTTPAddDocumentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddDocumentRequest

	if r.ContentLength == 0 {
		logger.Log("Add request with no body")
		return req, nil
	}
	//vars := mux.Vars(r)
	//team, ok := vars["team"]
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
	logger.Log("Inside encoder")
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
