package main

import (
	grpcclient "Book-repo/cmd/book/client"
	"Book-repo/internal"
	"flag"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	bookmark "Book-repo/pkg/book"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s\n\n%s\n", r, debug.Stack())
		}
	}()
	var (
		grpcAddr = flag.String("addr", ":8082",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()

	BookmarkService := grpcclient.New(conn)
	args := flag.Args()
	var cmd string
	cmd = pop(args)

	switch cmd {
	case "get":
		get(ctx, BookmarkService)
	case "servicestatus":
		servicestatus(ctx, BookmarkService)
	case "adddocument":
		addDocument(ctx, BookmarkService)
	case "Bookmark":
		Bookmark(ctx, BookmarkService)
	case "status":
		status(ctx, BookmarkService)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

//parse command line argument one by one
func pop(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

// call get service
func get(ctx context.Context, service bookmark.Service) {
	doc, err := service.Get(ctx, internal.Filter{Key: "Title", Value: "The Dark Code"})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(len(doc))
	for _, d := range doc {
		fmt.Println(d)
	}
}

func servicestatus(ctx context.Context, service bookmark.Service) {

	resp, err := service.ServiceStatus(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(resp)
}
func addDocument(ctx context.Context, service bookmark.Service) {
	ticketid, err := service.AddDocument(ctx, &internal.Document{Content: "Book",
		Title:  "The Dark Code",
		Author: "Bruce Wayne",
		Topic:  "Science",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Your ticketID is", ticketid)
}
func Bookmark(ctx context.Context, service bookmark.Service) {
	//newTicketID := shortuuid.New()
	resp, err := service.Bookmark(ctx,
		"pR6q6KabwpD6GvFd6PSBX5", 286)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if resp == 200 {
		fmt.Println("Updated succcessfully!!")
	} else {
		fmt.Println("Error processing your data..Please try again after sometime")
	}
}
func status(ctx context.Context, service bookmark.Service) {
	resp, err := service.Status(ctx, "adfgghg")
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(resp)
}
