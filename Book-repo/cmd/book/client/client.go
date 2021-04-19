package client

import (
	pb "Book-repo/api/v1/pb/book"
	Bookmark "Book-repo/pkg/book"
	"Book-repo/pkg/book/endpoints"
	"Book-repo/pkg/book/transport"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"google.golang.org/grpc"
)

// Return new Bookmark service
func New(conn *grpc.ClientConn) Bookmark.Service {

	var getEndpoint = grpctransport.NewClient(
		conn, "Bookmark.Bookmark", "Get",
		transport.EncodeGRPCGetRequest,
		transport.DecodeGRPCGetResponse,
		pb.GetReply{},
	).Endpoint()
	var addDocumentEndpoint = grpctransport.NewClient(
		conn, "Bookmark.Bookmark", "AddDocument",
		transport.EncodeGRPCAddDocumentRequest,
		transport.DecodeGRPCAddDocumentResponse,
		pb.AddDocumentReply{},
	).Endpoint()
	var statusEndpoint = grpctransport.NewClient(
		conn, "Bookmark.Bookmark", "Status",
		transport.EncodeGRPCStatusRequest,
		transport.DecodeGRPCStatusResponse,
		pb.StatusReply{},
	).Endpoint()
	var serviceStatusEndpoint = grpctransport.NewClient(
		conn, "Bookmark.Bookmark", "ServiceStatus",
		transport.EncodeGRPCServiceStatusRequest,
		transport.DecodeGRPCServiceStatusResponse,
		pb.ServiceStatusReply{},
	).Endpoint()
	var BookmarkEndpoint = grpctransport.NewClient(
		conn, "Bookmark.Bookmark", "Bookmark",
		transport.EncodeGRPCBookmarkRequest,
		transport.DecodeGRPCBookmarkResponse,
		pb.BookmarkReply{},
	).Endpoint()

	return &endpoints.Set{
		GetEndpoint:           getEndpoint,
		AddDocumentEndpoint:   addDocumentEndpoint,
		StatusEndpoint:        statusEndpoint,
		ServiceStatusEndpoint: serviceStatusEndpoint,
		BookmarkEndpoint:      BookmarkEndpoint,
	}
}
