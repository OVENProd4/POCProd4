package Bookmark

import (
	pb "Book-repo/api/v1/pb/db"
	"Book-repo/internal"
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/go-kit/kit/log"

	"google.golang.org/grpc"
)

type Service interface {
	Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error)
	Status(ctx context.Context, ticketID string) (internal.Status, error)
	Bookmark(ctx context.Context, ticketID string, mark int) (int, error)
	AddDocument(ctx context.Context, doc *internal.Document) (string, error)
	ServiceStatus(ctx context.Context) (int, error)
}

type BookmarkService struct {
	dbs pb.DatabaseClient
}

func NewService(conn *grpc.ClientConn) Service {
	var svc Service
	svc = NewBasicService(conn)
	return svc
}

//func NewService() Service { return &BookmarkService{} }

func NewBasicService(conn *grpc.ClientConn) Service {
	return &BookmarkService{
		dbs: pb.NewDatabaseClient(conn),
	}
}

func (w *BookmarkService) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	// query the database using the filters and return the list of documents
	// return error if the filter (key) is invalid and also return error if no item found
	logger.Log("Inside Get service")
	fil := make([]*pb.GetRequest_Filters, len(filters), len(filters))
	for i, f := range filters {

		fil[i] = &pb.GetRequest_Filters{
			Key:   f.Key,
			Value: f.Value,
		}
	}
	//	logger.Log(fil[0])
	resp, err := w.dbs.Get(ctx, &pb.GetRequest{Filters: fil})

	if err != nil {
		panic("Error calling db service")
	}
	docs := make([]internal.Document, len(resp.Documents), len(resp.Documents))
	for i, val := range resp.Documents {
		docs[i] = internal.Document{
			Content:  val.Content,
			Title:    val.Title,
			Author:   val.Author,
			Topic:    val.Topic,
			Bookmark: int(val.Bookmark),
		}
	}
	return docs, nil
}

func (w *BookmarkService) Status(_ context.Context, ticketID string) (internal.Status, error) {
	// query database using the ticketID and return the document info
	// return err if the ticketID is invalid or no Document exists for that ticketID
	return internal.InProgress, nil
}

func (w *BookmarkService) Bookmark(ctx context.Context, ticketID string, mark int) (int, error) {
	// update the database entry with Bookmark field as non empty
	// first check if the Bookmark status is not already in InProgress, Started or Finished state
	// If yes, then return invalid request
	// return error if no item found using the ticketID
	//var doc *pb.Document
	doc := &pb.Document{
		Bookmark: int64(mark),
	}
	//logger.Log(mark)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	resp, err := w.dbs.Update(ctx, &pb.UpdateRequest{TicketID: ticketID, Document: doc})
	if err != nil {
		panic("Error calling db service")
	} else {
		logger.Log("Success")
	}

	return int(resp.Code), nil
	//return http.StatusOK, nil
}

func (w *BookmarkService) AddDocument(ctx context.Context, doc *internal.Document) (string, error) {
	// add the document entry in the database by calling the database service
	// return error if the doc is invalid and/or the database invalid entry error
	logger.Log("Inside Bookmark service adddocument")

	document := &pb.Document{
		Author:  doc.Author,
		Content: doc.Content,
		Title:   doc.Title,
		Topic:   doc.Topic,
	}
	resp, err := w.dbs.Add(ctx, &pb.AddRequest{Document: document})
	if err != nil {
		panic("Error calling db service")
	}

	return resp.TicketID, nil
	//	return "aaaaa", nil
}

func (w *BookmarkService) ServiceStatus(ctx context.Context) (int, error) {
	logger.Log("Checking the Service health...")
	resp, err := w.dbs.ServiceStatus(ctx, &pb.ServiceStatusRequest{})
	if err != nil {
		panic("Error calling db service")
	}
	return int(resp.Code), nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
