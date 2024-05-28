package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"stats/db"
	pb "stats/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StatsServer struct {
	pb.UnimplementedStatisticsServiceServer
	db *sql.DB
}

func NewStatsServer() (*StatsServer, error) {
	db, err := db.NewDatabase(os.Getenv("CLICKHOUSE_URI"))
	if err != nil {
		return nil, err
	}
	return &StatsServer{db: db}, nil
}

func (s *StatsServer) GetTaskStats(ctx context.Context, req *pb.GetTaskStatsRequest) (*pb.GetTaskStatsResponse, error) {
	var likes uint64
	var views uint64

	s.db.QueryRow(`
		SELECT countDistinct(user_id)
		FROM likes
		WHERE task_id = $1
		GROUP BY task_id
	`, req.TaskId).Scan(&likes)

	s.db.QueryRow(`
		SELECT countDistinct(user_id)
		FROM views
		WHERE task_id = $1
		GROUP BY task_id
	`, req.TaskId).Scan(&views)

	return &pb.GetTaskStatsResponse{
		TaskId: req.TaskId,
		Likes:  likes,
		Views:  views,
	}, nil
}

func (s *StatsServer) GetTasksTop(ctx context.Context, req *pb.GetTasksTopRequest) (*pb.GetTasksTopResponse, error) {
	var table string
	switch *req.OrderBy.Enum() {
	case pb.StatsType_Likes:
		table = "likes"
	case pb.StatsType_Views:
		table = "views"
	}

	const numTasks = 5

	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT author_id, task_id, countDistinct(user_id) as count
		FROM %v
		GROUP BY task_id, author_id
		ORDER BY count DESC
		LIMIT %v
	`, table, numTasks),
	)

	if err == sql.ErrNoRows {
		return &pb.GetTasksTopResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []*pb.Task

	for rows.Next() {
		var authorID string
		var taskID string
		var count uint64
		if err := rows.Scan(&authorID, &taskID, &count); err != nil {
			return nil, err
		}
		tasks = append(tasks, &pb.Task{
			TaskId:   taskID,
			AuthorId: authorID,
			Count:    count,
		})
	}
	return &pb.GetTasksTopResponse{Tasks: tasks}, nil
}

func (s *StatsServer) GetUsersTop(ctx context.Context, req *emptypb.Empty) (*pb.GetUsersTopResponse, error) {
	const numUsers = 3

	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT author_id, COUNT(*) as total_likes
		FROM likes
		GROUP BY author_id
		ORDER BY total_likes DESC
		LIMIT %v
	`, numUsers),
	)

	if err == sql.ErrNoRows {
		return &pb.GetUsersTopResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	var users []*pb.User

	for rows.Next() {
		var authorID string
		var likes uint64
		if err := rows.Scan(&authorID, &likes); err != nil {
			return nil, err
		}
		users = append(users, &pb.User{
			UserId:     authorID,
			LikesCount: likes,
		})
	}
	return &pb.GetUsersTopResponse{Users: users}, nil
}
