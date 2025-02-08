package main

import (
	"Auth/service/Impl"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"project1/Auth/auth"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("AUTH_ADDRESS")+":"+os.Getenv("AUTH_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, Impl.NewAuthServiceImpl())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
