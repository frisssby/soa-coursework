package main

import (
	"main/db"
	"main/jwt"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func hashPassword(password string) string {
	// TODO
	return password
}

func signUp(c *gin.Context) {
	var userCreds models.UserCredentials
	if err := c.BindJSON(&userCreds); err != nil {
		return
	}

	_, err := db.GetUserCredentials(userCreds.Username)

	if err == nil {
		c.JSON(http.StatusForbidden, "User already exists")
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userCreds.Password = hashPassword(userCreds.Password)
	if err := db.CreateUser(userCreds); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, "Successfully signed up")
	}
}

func signIn(c *gin.Context) {
	var userCreds models.UserCredentials
	if err := c.BindJSON(&userCreds); err != nil {
		return
	}

	userCredsDB, err := db.GetUserCredentials(userCreds.Username)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusForbidden, "No user with provided username found")
		return
	} else if err != nil {
		panic(err)
	}

	if userCredsDB.Password != hashPassword(userCreds.Password) {
		c.JSON(http.StatusForbidden, "Wrong password")
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	tokenString, err := jwt.GenerateJWT(userCreds.Username, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Cannot generate token")
	}
	c.SetCookie("jwt", tokenString, expirationTime.Second(), "/", "localhost", false, true)
}

func updateUser(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "No token provided")
	}

	username, err := jwt.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	if username != c.Param("username") {
		c.JSON(http.StatusForbidden, "Not allowed to update this resource")
	}

	var userData models.UserData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	db.UpdateUserData(username, userData)
	c.JSON(http.StatusOK, "Successfully updated user data")
}
