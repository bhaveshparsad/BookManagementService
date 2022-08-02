package main

import (
	"fmt"
	"log"
	"net"

	pb "BookManagementService/protoFiles"

	c "BookManagementService/service"

	"google.golang.org/grpc"
)

const port = ":55001"

func main() {
	fmt.Println("Welcome to Book Management...")

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to Listen...", err)

	}

	//Server Intialisation
	server := grpc.NewServer()

	//Registering server as new grpc server
	pb.RegisterBookMgmtServiceServer(server, &c.BookMgmtServiceServer{})

	log.Println("Server Listening at", listen.Addr())

	if err := server.Serve(listen); err != nil {
		log.Fatal("Failed to Serve...", err)
	}
}