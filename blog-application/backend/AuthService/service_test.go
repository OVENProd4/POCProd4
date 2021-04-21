package main

import (
	"blog-application/global"
	"blog-application/proto"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "test@gmail.com", Username: "Carl", Password: string(pw)})
	server := authServer{}
	_, err := server.Login(context.Background(), &proto.LoginRequest{Login: "test@gmail.com", Password: "example"})
	if err != nil {
		t.Error("Error returned", err.Error())
	} else {
		t.Error("No Error")
	}
	// _, err := server.Login(context.Background(), &proto.LoginRequest{Login: "something", Password: "something"})
	// if err == nil {
	// 	t.Error("Error returned", err.Error())
	// }
}

func Test_authServer_UsernameUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Carl"})
	server := authServer{}
	_, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Carlo"})
	if err != nil {
		t.Error("Error returned", err.Error())
	} else {
		t.Error("Username already exists")
	}
}

func Test_authServer_EmailUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "test@gmail.com"})
	server := authServer{}
	_, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "test@yahoo.com"})
	if err != nil {
		t.Error("Error returned", err.Error())
	} else {
		t.Error("Email already exists")
	}
}

func Test_authServer_Signup(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "test@gmail.com", Username: "Carl123"})
	server := authServer{}
	_, err := server.Signup(context.Background(), &proto.SignupRequest{Email: "test@yoo.com", Username: "Carl13", Password: "test@12345"})
	if err != nil {
		t.Error("Error returned", err.Error())
	} else {
		t.Error("User Added")
	}
}

func Test_authServer_AuthUser(t *testing.T) {
	server := authServer{}
	res, err := server.AuthUser(context.Background(), &proto.AuthUserRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjYwN2QzODM4MTY5NzEwYThiOTA5ZTY2OVwiLFwiVXNlcm5hbWVcIjpcIkNhcmxcIixcIkVtYWlsXCI6XCJ0ZXN0QGdtYWlsLmNvbVwiLFwiUGFzc3dvcmRcIjpcIiQyYSQxMCRjbHFvOGdoL25CMXFQVzNwNW5EVzdlenpWbWhJZXJKZTJPVVVtRDAzQUs4cTc2SWxNbVVaZVwifSJ9.cFZJKpzFnpjaOrKR6zqxI5UUvxSQuqyqsiVMPcgH45s"})
	if err != nil {
		t.Error("an error returned")
	}
	if res.GetUsername() != "Carl" {
		t.Error("Wrong result returned")
	}
	if res.GetUsername() == "Carl" {
		t.Error("Carl returned")
	}

}
