package chat

import (
	"log"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
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

var users Users

// var user User

func readXML(cred string) bool {
	xmlFile, err := os.Open("users.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(byteValue, &users)

	for i := 0; i < len(users.Users); i++ {
		// fmt.Println("User Type: " + users.Users[i].Type)
		// fmt.Println("User Name: " + users.Users[i].Name)
		// fmt.Println("Password: " + users.Users[i].Password)

		if cred == users.Users[i].Name || cred == users.Users[i].Password {
			fmt.Println("User Name: " + users.Users[i].Name)
			fmt.Println("Password: " + users.Users[i].Password)
			return true
		}
	}

	return false
}

type Server struct {
}

var username string
var password string

func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)

	return &Message{
		Body: "Hello From the Server! \n---> Enter your Username: ",
		// Status: 1,
	}, nil
}

func (s *Server) EnterUserName(ctx context.Context, message *Message) (*Message, error) {

	if readXML(message.Body) {
		log.Printf("Username: %s", message.Body)
		username = message.Body
		return &Message{Body: "Enter your Password: "}, nil
	} else {
		log.Printf("Invalid Username: %s", message.Body)
		return &Message{Body: "Invalid Username, Enter your Username: "}, nil
	}
}

func (s *Server) EnterPassword(ctx context.Context, message *Message) (*Message, error) {
	if readXML(message.Body) {
		log.Printf("Password: %s", message.Body)
		password = message.Body
		return &Message{Body: "You have successfully Logged In as: " + username + " " + password}, nil
	} else {
		log.Printf("Invalid Password: %s", message.Body)
		return &Message{Body: "Invalid Password, Enter your Password: "}, nil
	}
}
