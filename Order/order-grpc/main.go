package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"order-grpc/order"
	"order-grpc/service/Impl"
	"os"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("ORDER_ADDRESS")+":"+os.Getenv("ORDER_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	order.RegisterOrderServiceServer(s, Impl.NewOrderService())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
