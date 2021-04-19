package transport

import (
	Bookmark "Book-repo/api/v1/pb/book"
	"context"
	"fmt"
	"runtime/debug"

	"Book-repo/pkg/book/endpoints"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpctransport.Handler
	status        grpctransport.Handler
	addDocument   grpctransport.Handler
	bookmark      grpctransport.Handler
	serviceStatus grpctransport.Handler
}

func NewGRPCServer(ep endpoints.Set) Bookmark.BookmarkServer {
	return &grpcServer{
		get: grpctransport.NewServer(
			ep.GetEndpoint,
			DecodeGRPCGetRequest,
			EncodeGRPCGetResponse,
		),
		status: grpctransport.NewServer(
			ep.StatusEndpoint,
			DecodeGRPCStatusRequest,
			EncodeGRPCStatusResponse,
		),
		addDocument: grpctransport.NewServer(
			ep.AddDocumentEndpoint,
			DecodeGRPCAddDocumentRequest,
			EncodeGRPCAddDocumentResponse,
		),
		bookmark: grpctransport.NewServer(
			ep.BookmarkEndpoint,
			DecodeGRPCBookmarkRequest,
			EncodeGRPCBookmarkResponse,
		),
		serviceStatus: grpctransport.NewServer(
			ep.ServiceStatusEndpoint,
			DecodeGRPCServiceStatusRequest,
			EncodeGRPCServiceStatusResponse,
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *Bookmark.GetRequest) (*Bookmark.GetReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*Bookmark.GetReply), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *Bookmark.ServiceStatusRequest) (*Bookmark.ServiceStatusReply, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	if err != nil {
		return nil, err
	}
	return rep.(*Bookmark.ServiceStatusReply), nil
}

func (g *grpcServer) AddDocument(ctx context.Context, r *Bookmark.AddDocumentRequest) (*Bookmark.AddDocumentReply, error) {
	_, rep, err := g.addDocument.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*Bookmark.AddDocumentReply), nil
}

func (g *grpcServer) Status(ctx context.Context, r *Bookmark.StatusRequest) (*Bookmark.StatusReply, error) {
	_, rep, err := g.status.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*Bookmark.StatusReply), nil
}

func (g *grpcServer) Bookmark(ctx context.Context, r *Bookmark.BookmarkRequest) (*Bookmark.BookmarkReply, error) {
	logger.Log("Inside Bookmark")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	_, rep, err := g.bookmark.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*Bookmark.BookmarkReply), nil
}
