package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"stats/event"
	"stats/db"
	pb "stats/proto"
	"stats/server"

	"google.golang.org/grpc"
)

func main() {
	db, err := db.NewDatabase(os.Getenv("CLICKHOUSE_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err.Error())
	}

	go event.RunConsumer("likes", db)
	go event.RunConsumer("views", db)

	statsServer, err := server.NewStatsServer(db)
	if err != nil {
		log.Fatalf("Failed to create statistics server: %v", err.Error())
	}

	grpcPortFlag := flag.Int("grpc-port", 51075, "grpc server port")
	httpPortFlag := flag.Int("http-port", 8080, "http server port")
	flag.Parse()

	go func() {
		log.Printf("Starting grpc server on port %d", *grpcPortFlag)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPortFlag))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterStatisticsServiceServer(grpcServer, statsServer)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve connections: ", err.Error())
		}
	}()

	http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPortFlag), nil); err != nil {
		log.Fatalf("Failed to serve http connections: %v", err.Error())
	}
}
