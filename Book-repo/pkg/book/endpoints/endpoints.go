package endpoints

import (
	"Book-repo/internal"
	Bookmark "Book-repo/pkg/book"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type GetRequest struct {
	Filters []internal.Filter `json:"filters,omitempty"`
}

type GetResponse struct {
	Documents []internal.Document `json:"documents"`
	Err       string              `json:"err,omitempty"`
}

type StatusRequest struct {
	TicketID string `json:"ticketID"`
}

type StatusResponse struct {
	Status internal.Status `json:"status"`
	Err    string          `json:"err,omitempty"`
}

type BookmarkRequest struct {
	TicketID   string `json:"ticketID"`
	PageNumber int    `json:"pagenumber"`
}

type BookmarkResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

type AddDocumentRequest struct {
	Document *internal.Document `json:"document"`
}

type AddDocumentResponse struct {
	TicketID string `json:"ticketID"`
	Err      string `json:"err,omitempty"`
}

type ServiceStatusRequest struct {
	Found bool `json:"found,omitempty"`
}

type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

type Set struct {
	GetEndpoint           endpoint.Endpoint
	AddDocumentEndpoint   endpoint.Endpoint
	StatusEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
	BookmarkEndpoint      endpoint.Endpoint
}

func NewEndpointSet(svc Bookmark.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(svc),
		AddDocumentEndpoint:   MakeAddDocumentEndpoint(svc),
		StatusEndpoint:        MakeStatusEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
		BookmarkEndpoint:      MakeBookmarkEndpoint(svc),
	}
}

func MakeGetEndpoint(svc Bookmark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		docs, err := svc.Get(ctx, req.Filters...)
		if err != nil {
			return GetResponse{docs, err.Error()}, nil
		}
		return GetResponse{docs, ""}, nil
	}
}

func MakeStatusEndpoint(svc Bookmark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		status, err := svc.Status(ctx, req.TicketID)
		if err != nil {
			return StatusResponse{Status: status, Err: err.Error()}, nil
		}
		return StatusResponse{Status: status, Err: ""}, nil
	}
}

func MakeAddDocumentEndpoint(svc Bookmark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddDocumentRequest)
		ticketID, err := svc.AddDocument(ctx, req.Document)
		if err != nil {
			return AddDocumentResponse{TicketID: ticketID, Err: err.Error()}, nil
		}
		return AddDocumentResponse{TicketID: ticketID, Err: ""}, nil
	}
}

func MakeBookmarkEndpoint(svc Bookmark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(BookmarkRequest)
		code, err := svc.Bookmark(ctx, req.TicketID, req.PageNumber)
		if err != nil {
			return BookmarkResponse{Code: code, Err: err.Error()}, nil
		}
		return BookmarkResponse{Code: code, Err: ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc Bookmark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	req := GetRequest{
		Filters: filters,
	}
	resp, err := s.GetEndpoint(ctx, req)
	if err != nil {
		return []internal.Document{}, err
	}
	getResp := resp.(GetResponse)
	if getResp.Err != "" {
		return []internal.Document{}, errors.New(getResp.Err)
	}
	return getResp.Documents, nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	req := ServiceStatusRequest{
		Found: true,
	}
	resp, err := s.ServiceStatusEndpoint(ctx, req)
	if err != nil {
		return -1, err
	}
	svcStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return svcStatusResp.Code, err
	}
	if svcStatusResp.Err != "" {
		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
	}
	return svcStatusResp.Code, nil
}

func (s *Set) AddDocument(ctx context.Context, doc *internal.Document) (string, error) {
	resp, err := s.AddDocumentEndpoint(ctx, AddDocumentRequest{Document: doc})
	if err != nil {
		return "", err
	}
	adResp := resp.(AddDocumentResponse)
	if adResp.Err != "" {
		return "", errors.New(adResp.Err)
	}
	return adResp.TicketID, nil
}

func (s *Set) Status(ctx context.Context, ticketID string) (internal.Status, error) {
	resp, err := s.StatusEndpoint(ctx, StatusRequest{TicketID: ticketID})
	if err != nil {
		return internal.Failed, err
	}
	stsResp := resp.(StatusResponse)
	if stsResp.Err != "" {
		return internal.Failed, errors.New(stsResp.Err)
	}
	return stsResp.Status, nil
}

func (s *Set) Bookmark(ctx context.Context, ticketID string, mark int) (int, error) {
	logger.Log("Inside Bookmark end")
	resp, err := s.BookmarkEndpoint(ctx, BookmarkRequest{TicketID: ticketID, PageNumber: mark})
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	wmResp := resp.(BookmarkResponse)
	if err != nil {
		return wmResp.Code, err
	}
	if wmResp.Err != "" {
		return wmResp.Code, errors.New(wmResp.Err)
	}
	return wmResp.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
