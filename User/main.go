package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"project1/User/service/Impl"
	"project1/User/user"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("USER_ADDRESS")+":"+os.Getenv("USER_PORT"))
	//listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, Impl.NewUserService())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
