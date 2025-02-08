package main

import (
	"cart-grpc/cart"
	"cart-grpc/service/Impl"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("CART_ADDRESS")+":"+os.Getenv("CART_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	cart.RegisterCartServiceServer(s, Impl.NewCartServiceImpl())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
