package handlers

import (
	"context"
	"net/http"
	"strconv"

	"users/event"
	"users/grpc"
	"users/models"
	statspb "users/proto/stats"
	taskpb "users/proto/tasks"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	username := c.GetString("username")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	response, err := grpc.TaskClient.CreateTask(context.Background(), &taskpb.CreateTaskRequest{
		UserId:      username,
		Description: task.Description,
		Status:      task.Status,
	})
	if err != nil {
		processTaskRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.Task{
		TaskID: response.TaskId,
	})
}

func UpdateTask(c *gin.Context) {
	username := c.GetString("username")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	taskID := c.Param("id")
	_, err := grpc.TaskClient.UpdateTask(context.Background(), &taskpb.UpdateTaskRequest{
		UserId:      username,
		TaskId:      taskID,
		Description: task.Description,
		Status:      task.Status,
	})
	if err != nil {
		processTaskRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, "Successfully updated task")
}

func DeleteTask(c *gin.Context) {
	username := c.GetString("username")
	taskID := c.Param("id")
	_, err := grpc.TaskClient.DeleteTask(context.Background(), &taskpb.DeleteTaskRequest{UserId: username, TaskId: taskID})
	if err != nil {
		processTaskRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, "Successfully deleted task")
}

func GetTask(c *gin.Context) {
	taskID := c.Param("id")
	resp, err := grpc.TaskClient.GetTask(context.Background(), &taskpb.GetTaskRequest{TaskId: taskID})
	if err != nil {
		processTaskRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.Task{
		UserID:      resp.AuthorId,
		TaskID:      resp.TaskId,
		Description: resp.Description,
		Status:      resp.Status,
	})
}

func GetTaskStats(c *gin.Context) {
	taskID := c.Param("id")
	response, err := grpc.StatsClient.GetTaskStats(context.Background(), &statspb.GetTaskStatsRequest{
		TaskId: taskID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "RPC failed")
		return
	}
	c.JSON(http.StatusOK, models.TaskStat{
		TaskID: taskID,
		Likes:  response.Likes,
		Views:  response.Views,
	})
}

func GetTopTasks(c *gin.Context) {
	orderBy, ok := c.GetQuery("order_by")
	if !ok {
		c.JSON(http.StatusBadRequest, "order_by query parameter required")
		return
	}
	var statsType statspb.StatsType
	switch orderBy {
	case "likes":
		statsType = statspb.StatsType_Likes
	case "views":
		statsType = statspb.StatsType_Views
	default:
		c.JSON(http.StatusBadRequest, "order_by query parameter should be either `likes` or `views`")
		return
	}
	response, err := grpc.StatsClient.GetTasksTop(context.Background(), &statspb.GetTasksTopRequest{
		OrderBy: statsType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "RPC failed")
		return
	}
	taskStats := make([]interface{}, 0)
	for _, stat := range response.GetTasks() {
		if orderBy == "likes" {
			taskStats = append(taskStats, models.LikesStat{
				TaskID: stat.TaskId,
				Likes:  stat.Count,
			})
		} else {
			taskStats = append(taskStats, models.TaskStat{
				TaskID: stat.TaskId,
				Views:  stat.Count,
			})
		}
	}
	c.JSON(http.StatusOK, taskStats)

}

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
		processTaskRPCError(c, err)
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
	resp, err := grpc.TaskClient.GetTask(context.Background(), &taskpb.GetTaskRequest{TaskId: taskID})
	if err != nil {
		return "", err
	}
	return resp.AuthorId, nil
}

const DEFAULT_PAGE_SIZE = 10

func getPageSize(c *gin.Context) int {
	pageSizeParam := c.Query("page_size")
	if pageSizeParam != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil {
			return pageSizeVal
		}
	}
	return DEFAULT_PAGE_SIZE
}

func getPageID(c *gin.Context) int {
	pageNumParam := c.Query("page_num")
	if pageNumParam != "" {
		if pageNumVal, err := strconv.Atoi(pageNumParam); err == nil {
			return max(1, pageNumVal) - 1
		}
	}
	return 0
}
