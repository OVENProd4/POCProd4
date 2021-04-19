package main

import (
	"POC/proto"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Server Running!!")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterAuthenticateServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) ValidateUser(ctx context.Context, req *proto.Request) (*proto.Response, error) {

	fmt.Println("Request received:", req)
	username := req.GetUsername()
	password := req.GetPassword()

	if username == "admin" && password == "123" {
		res := &proto.Response{
			Result: "Welcome User",
		}
		return res, nil
	} else {
		res := &proto.Response{
			Result: "Invalid Username/Password",
		}
		return res, nil
	}

}
