package transport

import (
	"context"
	"fmt"
	"runtime/debug"

	"Book-repo/api/v1/pb/db"
	"Book-repo/internal"
	"Book-repo/pkg/database/endpoints"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpctransport.Handler
	update        grpctransport.Handler
	add           grpctransport.Handler
	remove        grpctransport.Handler
	serviceStatus grpctransport.Handler
}

func NewGRPCServer(ep endpoints.Set) db.DatabaseServer {
	return &grpcServer{
		get: grpctransport.NewServer(
			ep.GetEndpoint,
			DecodeGRPCGetRequest,
			EncodeGRPCGetResponse,
		),
		update: grpctransport.NewServer(
			ep.UpdateEndpoint,
			DecodeGRPCUpdateRequest,
			EncodeGRPCUpdateResponse,
		),
		add: grpctransport.NewServer(
			ep.AddEndpoint,
			DecodeGRPCAddRequest,
			EncodeGRPCAddResponse,
		),
		remove: grpctransport.NewServer(
			ep.RemoveEndpoint,
			DecodeGRPRemoveRequest,
			EncodeGRPCRemoveResponse,
		),
		serviceStatus: grpctransport.NewServer(
			ep.ServiceStatusEndpoint,
			DecodeGRPCServiceStatusRequest,
			EncodeGRPCServiceStatusResponse,
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *db.GetRequest) (*db.GetReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*db.GetReply), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *db.ServiceStatusRequest) (*db.ServiceStatusReply, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*db.ServiceStatusReply), nil
}

func (g *grpcServer) Add(ctx context.Context, r *db.AddRequest) (*db.AddReply, error) {
	_, rep, err := g.add.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*db.AddReply), nil
}

func (g *grpcServer) Update(ctx context.Context, r *db.UpdateRequest) (*db.UpdateReply, error) {
	logger.Log("Inside db grpc")
	_, rep, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*db.UpdateReply), nil
}

func (g *grpcServer) Remove(ctx context.Context, r *db.RemoveRequest) (*db.RemoveReply, error) {
	_, rep, err := g.remove.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*db.RemoveReply), nil
}

func DecodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("DeodeGRPCGetRequest db")
	req := grpcReq.(*db.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func DecodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("Inside DecodeGRPCUpdateRequest db")
	req := grpcReq.(*db.UpdateRequest)
	doc := &internal.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int(req.Document.Bookmark),
	}
	return endpoints.UpdateRequest{TicketID: req.TicketID, Document: doc}, nil
}

func DecodeGRPRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.RemoveRequest)
	return endpoints.RemoveRequest{TicketID: req.TicketID}, nil
}

func DecodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.AddRequest)
	doc := &internal.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int(req.Document.Bookmark),
	}
	return endpoints.AddRequest{Document: doc}, nil
}

func DecodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func DecodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("DecodeGRPCGetResponse db")
	reply := grpcReply.(*db.GetReply)
	var docs []internal.Document
	for _, d := range reply.Documents {
		doc := internal.Document{
			Content:  d.Content,
			Title:    d.Title,
			Author:   d.Author,
			Topic:    d.Topic,
			Bookmark: int(d.Bookmark),
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs, Err: reply.Err}, nil
}

func DecodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("Inside DecodeGRPCUpdateReply db")
	reply := grpcReply.(*db.UpdateReply)
	return endpoints.UpdateResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

func DecodeGRPCRemoveResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*db.RemoveReply)
	return endpoints.RemoveResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

func DecodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*db.AddReply)
	return endpoints.AddResponse{TicketID: reply.TicketID, Err: reply.Err}, nil
}

func DecodeGRPCServiceStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*db.ServiceStatusReply)
	return endpoints.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}
func EncodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("EncodeGRPCGetRequest db")

	req := grpcReq.(endpoints.GetRequest)
	fil := make([]*db.GetRequest_Filters, len(req.Filters), len(req.Filters))
	for i, f := range req.Filters {

		fil[i] = &db.GetRequest_Filters{
			Key:   f.Key,
			Value: f.Value,
		}

	}

	return &db.GetRequest{Filters: fil}, nil
}

func EncodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("Inside EncodeGRPCUpdateRequest db")
	req := grpcReq.(endpoints.UpdateRequest)
	doc := &db.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int64(req.Document.Bookmark),
	}
	return db.UpdateRequest{TicketID: req.TicketID, Document: doc}, nil
}

func EncodeGRPRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(endpoints.RemoveRequest)
	return &db.RemoveRequest{TicketID: req.TicketID}, nil
}

func EncodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(endpoints.AddRequest)
	doc := &db.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int64(req.Document.Bookmark),
	}
	return &db.AddRequest{Document: doc}, nil
}

func EncodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return &db.ServiceStatusRequest{}, nil
}

func EncodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {

	logger.Log("EncodeGRPCGetResponse db")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()

	reply := grpcReply.(endpoints.GetResponse)
	var docs []*db.Document
	docs = make([]*db.Document, len(reply.Documents), len(reply.Documents))
	for i, d := range reply.Documents {
		doc := &db.Document{
			Content:  d.Content,
			Title:    d.Title,
			Author:   d.Author,
			Topic:    d.Topic,
			Bookmark: int64(d.Bookmark),
		}
		docs[i] = doc
	}
	return &db.GetReply{Documents: docs, Err: reply.Err}, nil
}

func EncodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("Inside EncodeGRPCUpdateRequest db")
	reply := grpcReply.(endpoints.UpdateResponse)
	return &db.UpdateReply{Code: int64(reply.Code), Err: reply.Err}, nil
}

func EncodeGRPCRemoveResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.RemoveResponse)
	return &db.RemoveReply{Code: int64(reply.Code), Err: reply.Err}, nil
}

func EncodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.AddResponse)
	return &db.AddReply{TicketID: reply.TicketID, Err: reply.Err}, nil
}

func EncodeGRPCServiceStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.ServiceStatusResponse)
	return &db.ServiceStatusReply{Code: int64(reply.Code), Err: reply.Err}, nil
}
