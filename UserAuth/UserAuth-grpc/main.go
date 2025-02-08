package main

import (
	"UserAuth-grpc/service/Impl"
	"UserAuth-grpc/userauth"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	godotenv.Load("config.env")
	listen, err := net.Listen("tcp", os.Getenv("USER_AUTH_ADDRESS")+":"+os.Getenv("USER_AUTH_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	userauth.RegisterAuthServiceServer(s, Impl.NewUserAuthServiceImpl())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
