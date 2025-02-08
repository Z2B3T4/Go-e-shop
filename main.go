package main

//
//import (
//	"google.golang.org/grpc"
//	"log"
//	"net"
//	"project1/User"
//)
//
//func main() {
//	listen, err := net.Listen("tcp", "localhost:9090")
//	if err != nil {
//		log.Fatal(err)
//	}
//	s := grpc.NewServer()
//
//	User.Userinit(s)
//
//	if err := s.Serve(listen); err != nil {
//		log.Fatal(err)
//	}
//	s.Stop()
//
//}
