package main

import (
	"POC/proto"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am Client")

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuthenticateServiceClient(conn)

	g := gin.Default()
	g.GET("/validate/:username/:password", func(c *gin.Context) {
		req := &proto.Request{
			Username: c.Param("username"),
			Password: c.Param("password"),
		}
		if resp, err := client.ValidateUser(c, req); err == nil {
			c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprint(resp.Result)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	if err := g.Run(":9090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
