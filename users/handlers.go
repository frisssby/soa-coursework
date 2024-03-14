package main

import (
	"users/db"
	"users/jwt"
	"users/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
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

	hashedPassword, err := hashPassword(userCreds.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	userCreds.Password = hashedPassword

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
	} 
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return;
	}

	if !checkPasswordHash(userCreds.Password, userCredsDB.Password) {
		c.JSON(http.StatusForbidden, "Wrong password")
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	tokenString, err := jwt.GenerateJWT(userCreds.Username, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Cannot generate token")
		return
	}
	c.SetCookie("jwt", tokenString, expirationTime.Second(), "/", "localhost", false, true)
}

func updateUser(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "No token provided")
		return
	}

	username, err := jwt.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	if username != c.Param("username") {
		c.JSON(http.StatusForbidden, "Not allowed to update this resource")
		return
	}

	var userData models.UserData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	db.UpdateUserData(username, userData)
	c.JSON(http.StatusOK, "Successfully updated user data")
}
