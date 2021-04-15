package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-tutorial/chat"

	"encoding/xml"
	"fmt"
)

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	XMLName  xml.Name `xml:"user"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name"`
	Password string   `xml:"password"`
}

// var users Users

// func readXML() {
// 	xmlFile, err := os.Open("users.xml")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Successfully Opened users.xml")
// 	defer xmlFile.Close()

// 	byteValue, _ := ioutil.ReadAll(xmlFile)

// 	xml.Unmarshal(byteValue, &users)

// 	// for i := 0; i < len(users.Users); i++ {
// 	// 	fmt.Println("User Type: " + users.Users[i].Type)
// 	// 	fmt.Println("User Name: " + users.Users[i].Name)
// 	// 	fmt.Println("Password: " + users.Users[i].Password)
// 	// }
// }

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	//////

	c := chat.NewChatServiceClient(conn)

	// readXML()

	//////

	message := chat.Message{
		Body: "Hello from the client!",
		// Status: 1,
	}

	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		// log.Fatalf("Error when calling SayHello: %s", err)
		fmt.Printf("Response from Server----- 1 ")

		fmt.Printf("Error when calling SayHello: %s", err)
	}

	fmt.Printf("Response from Server: %s", response.Body)

	/////
	fmt.Printf("Response from Server----- 2 ")
	var username string
	fmt.Printf("Response from Server----- 3 ")
	fmt.Scanf(" thagoo ", username)
	message1 := chat.Message{
		Body: username,
		// Status: 1,
	}
	fmt.Printf("Response from Server----- 4 ")
	response1, err1 := c.EnterUserName(context.Background(), &message1)
	if err1 != nil {
		log.Fatalf("Error when calling EnterUserName: %s", err1)
	}

	fmt.Printf("---> %s", response1.Body)

	//////
	var password string
	fmt.Scanln(&password)
	message2 := chat.Message{
		Body: password,
		// Status: 1,
	}

	response2, err2 := c.EnterPassword(context.Background(), &message2)
	if err2 != nil {
		log.Fatalf("Error when calling EnterPassword: %s", err2)
	}

	fmt.Printf("---> %s", response2.Body)

}
