package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "stats/proto"
	"stats/server"

	"google.golang.org/grpc"
)

func main() {
	statsServer, err := server.NewStatsServer();
	if (err != nil) {
		log.Fatalf("Failed to create statistics server: %v", err.Error())
	}

	portFlag := flag.Int("port", 51075, "grpc server port")
	flag.Parse()

	log.Printf("Starting grpc server on port %d", *portFlag)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatisticsServiceServer(grpcServer, statsServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve connections: ", err.Error())
	}
}
