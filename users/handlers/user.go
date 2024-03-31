package handlers

import (
	"log"
	"net/http"

	"users/db"
	"users/models"

	"github.com/gin-gonic/gin"
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
