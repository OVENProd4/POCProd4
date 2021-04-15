package chat

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

var username string
var password string

func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &Message{Body: "Hello From the Server! \nEnter your Username: "}, nil
}

func (s *Server) EnterUserName(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Userbame: %s", message.Body)
	username = message.Body
	return &Message{Body: "Enter your Password: "}, nil
}
func (s *Server) EnterPassword(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Password: %s", message.Body)
	password = message.Body
	return &Message{Body: "You have successfully Logged In as: " + username + " " + password}, nil
}