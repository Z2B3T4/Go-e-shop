package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"paymentGrpc/payment"
	"paymentGrpc/service/Impl"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("PAYMENT_ADDRESS")+":"+os.Getenv("PAYMENT_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	payment.RegisterPaymentServiceServer(s, Impl.NewPaymentService())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
