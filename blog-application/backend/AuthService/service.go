package main

import (
	"blog-application/global"
	"blog-application/proto"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type authServer struct{}

func (as authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	login := in.GetLogin()
	password := in.GetPassword()

	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{{"username": login}, {"email": login}}}).Decode(&user)
	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("wrong Login Credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("wrong Login Credentials")
	}
	return &proto.AuthResponse{Token: user.GetToken()}, nil
}

func (server authServer) Signup(_ context.Context, in *proto.SignupRequest) (*proto.AuthResponse, error) {
	username := in.GetUsername()
	email := in.GetEmail()
	password := in.GetPassword()

	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(username) < 4 || len(username) > 20 || len(email) < 7 || len(email) > 35 || len(password) < 8 || len(username) > 20 || !match {
		return &proto.AuthResponse{}, errors.New("validation failed")
	}
	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Println("Error returned from UsernameUsed:", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}
	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("username already exists")
	}

	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: email})
	if err != nil {
		log.Println("Error returned from EmailUsed:", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}
	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("email already exists")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := global.User{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw)}
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = global.DB.Collection("user").InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error inserting new user", err.Error())
		return &proto.AuthResponse{}, errors.New("something went wrong")
	}
	return &proto.AuthResponse{Token: newUser.GetToken()}, nil
}
func (as authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	email := in.GetEmail()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}
func (as authServer) UsernameUsed(_ context.Context, in *proto.UsernameUsedRequest) (*proto.UsedResponse, error) {
	username := in.GetUsername()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}
func (as authServer) AuthUser(_ context.Context, in *proto.AuthUserRequest) (*proto.AuthUserResponse, error) {
	token := in.GetToken()
	user := global.UserFromToken(token)
	return &proto.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil

}
func main() {
	fmt.Println("Server started!!")
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Error creating listener:%v", err)
	}
	go func() {
		log.Fatal("Serving gRPC: ", server.Serve(listener).Error())
	}()

	grpcWebServer := grpcweb.WrapServer(server)

	httpServer := &http.Server{
		Addr: ":9001",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())
}
