package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	gr "Orders/grpc_server"
	pb "Orders/proto"
	rt "Orders/rest_server"
)

const (
	grpcPort = "50051" //Server
	restPort = "8080"  //Client
)

// app is a convenience wrapper for all things needed to start
// and shutdown the Order microservice
type app struct {
	restServer rt.RestServer
	grpcServer gr.GrpcServer
	/* Listens for an application termination signal
	   Ex. (Ctrl X, Docker container shutdown, etc) */
	shutdownCh chan os.Signal
}

type Status int32

const (
	PENDING   = 0
	PAID      = 1
	SHIPPED   = 2
	DELIVERED = 3
	CANCELLED = 4
)

type msg struct {
	order_id int64
	//	items    Item
	total  float32
	status Status
}

//8080-client
//50051- server

// start starts the REST and gRPC Servers in the background
func (a app) start() {
	go a.restServer.Start() // non blocking now
	a.grpcServer.Start()    // also non blocking :-)
}

// stop shuts down the servers
func (a app) shutdown() error {
	a.grpcServer.Stop()
	return a.restServer.Stop()
}

type server struct {
	pb.UnimplementedOrderServiceServer
}

// newApp creates a new app with REST & gRPC servers
// this func performs all app related initialization
func newApp() (app, error) {
	//orderService := pb.UnimplementedOrderServiceServer{}
	//orderSevice :=
	gs, err := gr.NewGrpcServer(&server{}, grpcPort)
	if err != nil {
		return app{}, err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	return app{
		restServer: rt.NewRestServer(&server{}, restPort),
		grpcServer: gs,
		shutdownCh: quit,
	}, nil
}

func (c *server) Retrieve(ctx context.Context, request *pb.RetrieveOrderRequest) (*pb.RetrieveOrderResponse, error) {
	fmt.Println("Order Status")
	order_no := request.GetOrderId()
	fmt.Println("Order Status", order_no)
	mapTest := make(map[int64]msg)
	mapTest[10] = msg{
		order_id: 1,
		total:    PENDING,
		status:   1}
	mapTest[11] = msg{
		order_id: 2,
		total:    2,
		status:   1}
	mapTest[12] = msg{
		order_id: 3,
		total:    1,
		status:   3}
	mapTest[13] = msg{
		order_id: 4,
		total:    1,
		status:   2}
	fmt.Println(mapTest)
	order_info := new(pb.Order)
	keys := reflect.ValueOf(mapTest).MapKeys()
	fmt.Println(keys)
	for key, value := range mapTest {
		fmt.Println("Existing Order 1111111111111111111111111111111!!!!!!!!!!!!", key, value)
		if order_no == int64(key) {
			fmt.Println("Existing Order 222222222222222222222222222!!!!!!!!!!!!")
			order_info.OrderId = mapTest[order_no].order_id
			order_info.Status = pb.Order_Status(mapTest[order_no].status)
			order_info.Total = mapTest[order_no].total
			fmt.Println(order_info)
		} else {
			//return &pb.RetrieveOrderResponse{Order: nil}, nil
			order_info.Status = pb.Order_Status(-1)
			fmt.Println("Not Existing Order!!!!!!!!!!!!")
		}
	}

	return &pb.RetrieveOrderResponse{Order: order_info}, nil
}

func (c *server) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	fmt.Println("Order Status")
	order_no := request.GetOrderId()
	fmt.Println("Order Status", order_no)
	order_info := new(pb.Order)
	if order_no == 10 {
		order_info.OrderId = 500
		order_info.Status = 2
		order_info.Total = 4
	} else if order_no == 11 {
		order_info.OrderId = 11
		order_info.Status = 1
		order_info.Total = 4
	} else {
		return &pb.UpdateOrderResponse{Order: nil}, nil
	}
	return &pb.UpdateOrderResponse{Order: order_info}, nil
}

// run starts the app, handling any REST or gRPC server error
// and any shutdown signal
func run() error {
	app, err := newApp()
	if err != nil {
		return err
	}

	app.start()
	defer app.shutdown()

	select {
	case restErr := <-app.restServer.Error():
		return restErr
	case grpcErr := <-app.grpcServer.Error():
		return grpcErr
	case <-app.shutdownCh:
		return nil
	}
}

func main() {
	fmt.Println("In main")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
