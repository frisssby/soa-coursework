package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"tasks/db"
	pb "tasks/proto/tasks"
	"tasks/server"

	"google.golang.org/grpc"
)

func main() {
	if err := db.ConnectToMongoDB(os.Getenv("MONGODB_URI")); err != nil {
		log.Fatal("Could not connect to MongoDB", err.Error())
	}

	portFlag := flag.Int("port", 51075, "http server port")
	flag.Parse()

	log.Printf("Starting grpc server on port %d", *portFlag)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTaskServiceServer(grpcServer, &server.TaskServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve connections: ", err.Error())
	}
}
