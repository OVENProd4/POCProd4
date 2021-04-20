package Orders

import (
	pb "Orders/proto"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer implements a gRPC Server for the Order service
type GrpcServer struct {
	server   *grpc.Server
	errCh    chan error
	listener net.Listener
}
type server struct {
	pb.UnimplementedOrderServiceServer
}

// NewGrpcServer is a convenience func to create a GrpcServer
func NewGrpcServer(service pb.OrderServiceServer, port string) (GrpcServer, error) {
	fmt.Println("Hey man its inside NewGrpcServer!!!!!, port", port)

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, service)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return GrpcServer{}, err
	}

	return GrpcServer{
		server:   server,
		listener: lis,
		errCh:    make(chan error),
	}, nil
}

// Start starts the server in the background, pushing any error to the error channel
func (g GrpcServer) Start() {
	go func() {
		if err := g.server.Serve(g.listener); err != nil {
			fmt.Println("Hey Man its grpc!!!!!!!")
			g.errCh <- err
		}
	}()
}

// Stop stops the gRPC server
func (g GrpcServer) Stop() {
	g.server.GracefulStop()
}

// Error returns the server's error channel
func (g GrpcServer) Error() chan error {
	return g.errCh
}
