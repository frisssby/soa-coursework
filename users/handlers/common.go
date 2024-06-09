package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
