package grpc

import (
	pb "users/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TaskClient pb.TaskServiceClient

func InitTasksClient(uri string) error {
	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	TaskClient = pb.NewTaskServiceClient(conn)
	return nil
}
