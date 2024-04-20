package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"users/grpc"
	"users/models"
	pb "users/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateTask(c *gin.Context) {
	username := c.GetString("username")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	response, err := grpc.TaskClient.CreateTask(context.Background(), &pb.CreateTaskRequest{
		UserId:      username,
		Description: task.Description,
		Status:      task.Status,
	})
	if err != nil {
		processRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.Task{
		TaskId: response.TaskId,
	})
}

func UpdateTask(c *gin.Context) {
	username := c.GetString("username")
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	taskId := c.Param("id")
	_, err := grpc.TaskClient.UpdateTask(context.Background(), &pb.UpdateTaskRequest{
		UserId:      username,
		TaskId:      taskId,
		Description: task.Description,
		Status:      task.Status,
	})
	if err != nil {
		processRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, "Successfully updated task")
}

func DeleteTask(c *gin.Context) {
	username := c.GetString("username")
	taskId := c.Param("id")
	_, err := grpc.TaskClient.DeleteTask(context.Background(), &pb.DeleteTaskRequest{UserId: username, TaskId: taskId})
	if err != nil {
		processRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, "Successfully deleted task")
}

func ListTasks(c *gin.Context) {
	username := c.GetString("username")
	log.Println("page_size", getPageSize(c), "id", getPageID(c))
	resp, err := grpc.TaskClient.ListTasks(context.Background(), &pb.ListTasksRequest{
		UserId:   username,
		PageId:   int32(getPageID(c)),
		PageSize: int32(getPageSize(c)),
	})

	if err != nil {
		processRPCError(c, err)
		return
	}
	tasks := make([]models.Task, 0)
	for _, task := range resp.Tasks {
		tasks = append(tasks, models.Task{
			TaskId:      task.TaskId,
			Description: task.Description,
			Status:      task.Status,
		})
	}
	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	username := c.GetString("username")
	taskId := c.Param("id")
	resp, err := grpc.TaskClient.GetTask(context.Background(), &pb.GetTaskRequest{UserId: username, TaskId: taskId})
	if err != nil {
		processRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.Task{Description: resp.Description, Status: resp.Status})
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

func processRPCError(c *gin.Context, err error) {
	if status.Code(err) == codes.NotFound {
		c.JSON(http.StatusNotFound, "Task not found")
		return
	}
	if status.Code(err) == codes.PermissionDenied {
		c.JSON(http.StatusForbidden, "Not enough rights")
		return
	}
	log.Println("RPC failed", err.Error())
	c.JSON(http.StatusInternalServerError, "RPC failed")
}
