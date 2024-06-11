package grpc

import (
	pb "users/proto/stats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var StatsClient pb.StatisticsServiceClient

func InitStatsClient(uri string) error {
	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	StatsClient = pb.NewStatisticsServiceClient(conn)
	return nil
}
