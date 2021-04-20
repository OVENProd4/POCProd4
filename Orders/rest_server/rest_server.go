package Orders

import (
	pb "Orders/proto"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RestServer implements a REST server for the Order Service.
type RestServer struct {
	server       *http.Server
	orderService pb.OrderServiceServer // the same order service we injected into the gRPC server
	errCh        chan error
}

// NewRestServer is a convenience func to create a RestServer
func NewRestServer(orderService pb.OrderServiceServer, port string) RestServer {
	router := gin.Default()

	rs := RestServer{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		orderService: orderService,
		errCh:        make(chan error),
	}
	fmt.Println("port", port)

	//m := make(map[int]pb.Order)

	// register routes
	//--GET---------------------------------------------------------------
	router.GET("/order/:id", func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
			return
		}
		req := &pb.RetrieveOrderRequest{OrderId: int64(id)}
		fmt.Println("Hey Man its inside Get 1111111111111")
		if response, err := rs.orderService.Retrieve(ctx, req); err == nil {
			{
				if response.Order.Status != -1 {
					ctx.JSON(http.StatusOK, gin.H{
						"Order": fmt.Sprint(response.Order)})
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"Order": "Details are not found"})
				}
			}
		} else {
			//fmt.Println("Hey Man its inside Get 333333333333333333")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	//PUT-------------------------------------------------------------
	router.PUT("/order/1/:id", func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
			return
		}
		req := &pb.UpdateOrderRequest{OrderId: int64(id)}
		//fmt.Println("Hey Man its inside Get 1111111111111")
		if response, err := rs.orderService.Update(ctx, req); err == nil {
			fmt.Println("Hey Man its inside Get 22222222222222")
			ctx.JSON(http.StatusOK, gin.H{
				"Order": fmt.Sprint(response.Order),
			})
		} else {
			fmt.Println("Hey Man its inside Get 333333333333333333")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	router.DELETE("/order", rs.delete)
	router.GET("/order", rs.list)

	return rs
}

// Start starts the REST server in the background, pushing an error to the error channel
func (r RestServer) Start() {
	go func() {
		fmt.Println("Hey Man its rest!!!!!!!!!!")
		if err := r.server.ListenAndServe(); err != nil {
			r.errCh <- err
		}

	}()

}

// Stop stops the server
func (r RestServer) Stop() error {
	return r.server.Close()
}

// Error returns the server's error channel
func (r RestServer) Error() chan error {
	return r.errCh
}

// create is a handler func that creates an order from an order request (JSON body

/*func (r RestServer) retrieve(c *gin.Context) {
	name := ctx.Param("id")
	var req pb.RetrieveOrderRequest
	// use the order service to create the order from the req
	resp, err := r.orderService.Retrieve(c.Request.Context(), &req)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating order")
	}
	m := &jsonpb.Marshaler{}
	if err := m.Marshal(c.Writer, resp); err != nil {
		c.String(http.StatusInternalServerError, "error sending order response")
	}
	c.String(http.StatusNotImplemented, "not implemented yet")

}*/

func (r RestServer) delete(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}

func (r RestServer) list(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented yet")
}
