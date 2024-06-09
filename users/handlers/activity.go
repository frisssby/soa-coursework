package handlers

import (
	"context"
	"net/http"

	"users/event"
	"users/grpc"
	"users/models"
	taskspb "users/proto/tasks"

	"github.com/gin-gonic/gin"
)

func LikeTask(c *gin.Context) {
	handleActivity(c, "likes")
}

func ViewTask(c *gin.Context) {
	handleActivity(c, "views")
}

func handleActivity(c *gin.Context, brokerTopic string) {
	username := c.GetString("username")
	taskID := c.Param("id")
	authorID, err := getTaskAuthor(taskID)
	if err != nil {
		processRPCError(c, err)
		return
	}
	if err := event.ProduceEvent(brokerTopic, models.Event{
		UserID:   username,
		TaskID:   taskID,
		AuthorID: authorID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to write message to broker")
	}
}

func getTaskAuthor(taskID string) (string, error) {
	resp, err := grpc.TaskClient.GetTask(context.Background(), &taskspb.GetTaskRequest{TaskId: taskID})
	if err != nil {
		return "", err
	}
	return resp.AuthorId, nil
}
