package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "Book-repo/api/v1/pb/book"
	Bookmark "Book-repo/pkg/book"
	"Book-repo/pkg/book/endpoints"
	"Book-repo/pkg/book/transport"

	"google.golang.org/grpc"
)

func main() {

	var (
		gRPCAddr = flag.String("grpc", ":8081",
			"gRPC listen address")
	)
	flag.Parse()
	//ctx := context.Background()

	// init lorem service
	var svc Bookmark.Service
	svc = Bookmark.NewService()
	errChan := make(chan error)

	// creating Endpoints struct
	endpoints := endpoints.NewEndpointSet(svc)

	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			errChan <- err
			return
		}

		handler := transport.NewGRPCServer(endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterBookmarkServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
