package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"payment-grpc/checkout"
	"payment-grpc/service/Impl"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("CHECKOUT_ADDRESS")+":"+os.Getenv("CHECKOUT_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	checkout.RegisterCheckoutServiceServer(s, Impl.NewPaymentService())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
