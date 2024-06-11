package handlers

import (
	"context"
	"log"
	"net/http"

	"users/db"
	"users/grpc"
	"users/models"
	taskspb "users/proto/tasks"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

func UpdateUser(c *gin.Context) {
	username := c.GetString("username")
	if username != c.Param("username") {
		c.JSON(http.StatusForbidden, "Not allowed to update this resource")
		return
	}

	var userData models.UserData
	if err := c.BindJSON(&userData); err != nil {
		return
	}

	if err := db.UpdateUserData(username, userData); err != nil {
		log.Println("Error updating user data: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusOK, "Successfully updated user data")
}

func ListTasks(c *gin.Context) {
	userID := c.Param("id")
	resp, err := grpc.TaskClient.ListTasks(context.Background(), &taskspb.ListTasksRequest{
		AuthorId: userID,
		PageId:   int32(getPageID(c)),
		PageSize: int32(getPageSize(c)),
	})
	if err != nil {
		processTaskRPCError(c, err)
		return
	}
	tasks := make([]models.Task, 0)
	for _, task := range resp.Tasks {
		tasks = append(tasks, models.Task{
			TaskID:      task.TaskId,
			Description: task.Description,
			Status:      task.Status,
		})
	}
	c.JSON(http.StatusOK, tasks)
}

func GetTopUsers(c *gin.Context) {
	response, err := grpc.StatsClient.GetUsersTop(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "RPC failed")
		return
	}
	userStats := make([]models.UserStat, 0)
	for _, stat := range response.GetUsers() {
		userStats = append(userStats, models.UserStat{
			UserID: stat.UserId,
			Likes:  stat.LikesCount,
		})
	}
	c.JSON(http.StatusOK, userStats)
}
