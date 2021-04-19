package transport

import (
	Bookmark "Book-repo/api/v1/pb/book"
	"Book-repo/internal"
	"Book-repo/pkg/book/endpoints"
	"context"
	"fmt"
	"runtime/debug"
)

//Decode Get Request
func DecodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*Bookmark.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

//Decode Status Request
func DecodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*Bookmark.StatusRequest)
	return endpoints.StatusRequest{TicketID: req.TicketID}, nil
}

//Decode Bookmark Request
func DecodeGRPCBookmarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("Inside decodeGRPCBookmarkReq")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	req := grpcReq.(*Bookmark.BookmarkRequest)
	return endpoints.BookmarkRequest{TicketID: req.TicketID, PageNumber: int(req.Pagenumber)}, nil
}

//Decode AddDocument Request
func DecodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*Bookmark.AddDocumentRequest)
	doc := &internal.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int(req.Document.Bookmark),
	}
	return endpoints.AddDocumentRequest{Document: doc}, nil
}

//Decode ServiceStatus Request
func DecodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*Bookmark.ServiceStatusRequest)
	logger.Log("DecodeGRPCServiceStatusRequest")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	return endpoints.ServiceStatusRequest{Found: req.Found}, nil
}

//Decode Get Response
func DecodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("DecodeGRPCGetResponse")
	reply := grpcReply.(*Bookmark.GetReply)
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

//Encode and Decode Status Response
func DecodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*Bookmark.StatusReply)
	if reply.Status == Bookmark.StatusReply_STARTED {
		return endpoints.StatusResponse{Status: internal.Started, Err: reply.Err}, nil
	} else if reply.Status == Bookmark.StatusReply_IN_PROGRESS {
		return endpoints.StatusResponse{Status: internal.InProgress, Err: reply.Err}, nil
	} else if reply.Status == Bookmark.StatusReply_PENDING {
		return endpoints.StatusResponse{Status: internal.Pending, Err: reply.Err}, nil
	} else if reply.Status == Bookmark.StatusReply_FINISHED {
		return endpoints.StatusResponse{Status: internal.Finished, Err: reply.Err}, nil
	}
	return endpoints.StatusResponse{Status: internal.Failed, Err: reply.Err}, nil
}

//Encode and Decode Bookmark Response
func DecodeGRPCBookmarkResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("Inside decodeGRPCBookmarkResp")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	reply := grpcReply.(*Bookmark.BookmarkReply)
	return endpoints.BookmarkResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

//Encode and Decode AddDocument Response
func DecodeGRPCAddDocumentResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*Bookmark.AddDocumentReply)
	return endpoints.AddDocumentResponse{TicketID: reply.TicketID, Err: reply.Err}, nil
}

//Encode and Decode ServiceStatus Response
func DecodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*Bookmark.ServiceStatusReply)
	return endpoints.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}
func EncodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("EncodeGRPCGetRequest")
	req := grpcReq.(endpoints.GetRequest)
	//var fil []*Bookmark.GetRequest_Filters
	fil := make([]*Bookmark.GetRequest_Filters, len(req.Filters), len(req.Filters))
	for i, f := range req.Filters {

		fil[i] = &Bookmark.GetRequest_Filters{
			Key:   f.Key,
			Value: f.Value,
		}

	}

	return &Bookmark.GetRequest{Filters: fil}, nil
}

//Encode Status Request
func EncodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(endpoints.StatusRequest)
	return &Bookmark.StatusRequest{TicketID: req.TicketID}, nil
}

//Encode Bookmark Request
func EncodeGRPCBookmarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	logger.Log("Inside encodeGRPCBookmark")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	req := grpcReq.(endpoints.BookmarkRequest)
	return &Bookmark.BookmarkRequest{TicketID: req.TicketID, Pagenumber: int64(req.PageNumber)}, nil
}

//Encode AddDocument Request
func EncodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(endpoints.AddDocumentRequest)
	doc := &Bookmark.Document{
		Content:  req.Document.Content,
		Title:    req.Document.Title,
		Author:   req.Document.Author,
		Topic:    req.Document.Topic,
		Bookmark: int64(req.Document.Bookmark),
	}
	return &Bookmark.AddDocumentRequest{Document: doc}, nil
}

//Encode ServiceStatus Request
func EncodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	empty := grpcReq.(endpoints.ServiceStatusRequest)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	return &Bookmark.ServiceStatusRequest{Found: empty.Found}, nil
}

//Encode Get Response
func EncodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.GetResponse)
	var docs []*Bookmark.Document
	docs = make([]*Bookmark.Document, len(reply.Documents), len(reply.Documents))
	for i, d := range reply.Documents {
		doc := &Bookmark.Document{
			Content:  d.Content,
			Title:    d.Title,
			Author:   d.Author,
			Topic:    d.Topic,
			Bookmark: int64(d.Bookmark),
		}
		docs[i] = doc
	}
	return &Bookmark.GetReply{Documents: docs, Err: reply.Err}, nil
}

//Encode and Encode Status Response
func EncodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.StatusResponse)
	if reply.Status == internal.Started {
		return &Bookmark.StatusReply{Status: Bookmark.StatusReply_STARTED, Err: reply.Err}, nil
	} else if reply.Status == internal.InProgress {
		return &Bookmark.StatusReply{Status: Bookmark.StatusReply_IN_PROGRESS, Err: reply.Err}, nil
	} else if reply.Status == internal.Pending {
		return &Bookmark.StatusReply{Status: Bookmark.StatusReply_PENDING, Err: reply.Err}, nil
	} else if reply.Status == internal.Finished {
		return &Bookmark.StatusReply{Status: Bookmark.StatusReply_FINISHED, Err: reply.Err}, nil
	}
	return &Bookmark.StatusReply{Status: Bookmark.StatusReply_FAILED, Err: reply.Err}, nil
}

//Encode and Encode Bookmark Response
func EncodeGRPCBookmarkResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	logger.Log("Inside encodeGRPCBookmarkResp")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	reply := grpcReply.(endpoints.BookmarkResponse)
	return &Bookmark.BookmarkReply{Code: int64(reply.Code), Err: reply.Err}, nil
}

//Encode and Encode AddDocument Response
func EncodeGRPCAddDocumentResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.AddDocumentResponse)
	return &Bookmark.AddDocumentReply{TicketID: reply.TicketID, Err: reply.Err}, nil
}

//Encode and Encode ServiceStatus Response
func EncodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(endpoints.ServiceStatusResponse)
	if reply.Code != -1 {
		return &Bookmark.ServiceStatusReply{Code: int64(reply.Code), Err: reply.Err}, nil
	}
	return &Bookmark.ServiceStatusReply{Code: -1, Err: reply.Err}, nil
}
