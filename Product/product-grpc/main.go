package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"product-grpc/product"
	Impl2 "product-grpc/service/Impl"
)

func main() {
	godotenv.Load("./config.env")
	listen, err := net.Listen("tcp", os.Getenv("PRODUCT_ADDRESS")+":"+os.Getenv("PRODUCT_PORT"))
	//listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	//user.RegisterUserServiceServer(s, Impl.NewUserService())
	product.RegisterProductCatalogServiceServer(s, Impl2.NewProductService())
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	s.Stop()

}
