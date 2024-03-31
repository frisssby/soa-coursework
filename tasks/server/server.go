package server

import (
	"context"
	"log"

	"tasks/db"
	"tasks/models"

	pb "tasks/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	newTask, err := db.CreateTask(&models.Task{
		UserId:      req.UserId,
		Description: req.Content,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := pb.CreateTaskResponse{
		TaskId: newTask.TaskId,
	}
	return &resp, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*emptypb.Empty, error) {
	task := models.Task{
		UserId:      req.UserId,
		Description: req.Content,
	}
	stored, err := db.GetTask(req.TaskId)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if stored.UserId != req.UserId {
		log.Println("stored", stored.UserId, "req", req.UserId)
		return nil, status.Error(codes.PermissionDenied, "not enough rights")
	}
	if err := db.UpdateTask(req.TaskId, task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*emptypb.Empty, error) {
	stored, err := db.GetTask(req.TaskId)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if stored.UserId != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "not enough rights")
	}
	if err := db.DeleteTask(req.TaskId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	stored, err := db.GetTask(req.TaskId)
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if stored.UserId != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "not enough rights")
	}
	resp := pb.GetTaskResponse{
		TaskId:  stored.TaskId,
		Content: stored.Description,
	}
	return &resp, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := db.ListsTasks(req.UserId, int(req.PageSize), int(req.PageId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var resp pb.ListTasksResponse
	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, &pb.GetTaskResponse{TaskId: task.TaskId, Content: task.Description})
	}
	return &resp, nil
}
