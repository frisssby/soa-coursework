package handlers

import (
	"log"
	"net/http"
	"time"

	"users/db"
	"users/jwt"
	"users/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTl = 60 * time.Minute

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignUp(c *gin.Context) {
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
		log.Println("Error finding user creadential: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	hashedPassword, err := hashPassword(userCreds.Password)
	if err != nil {
		log.Println("Failed to hash password: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	userCreds.Password = hashedPassword

	if err := db.CreateUser(userCreds); err != nil {
		log.Println("Error creating user: ", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(tokenTTl)
	tokenString, err := jwt.GenerateJWT(userCreds.Username, expirationTime)
	if err != nil {
		log.Println("Error generating token: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.SetCookie("jwt", tokenString, expirationTime.Second(), "/", "localhost", false, true)
	c.JSON(http.StatusOK, "Successfully signed up")
}

func SignIn(c *gin.Context) {
	var userCreds models.UserCredentials
	if err := c.BindJSON(&userCreds); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	userCredsDB, err := db.GetUserCredentials(userCreds.Username)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusForbidden, "No user with provided username found")
		return
	}
	if err != nil {
		log.Println("Error finding user creadential: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	if !checkPasswordHash(userCreds.Password, userCredsDB.Password) {
		c.JSON(http.StatusForbidden, "Wrong password")
		return
	}

	expirationTime := time.Now().Add(tokenTTl)
	tokenString, err := jwt.GenerateJWT(userCreds.Username, expirationTime)
	if err != nil {
		log.Println("Error generating token: ", err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.SetCookie("jwt", tokenString, expirationTime.Second(), "/", "localhost", false, true)
	c.JSON(http.StatusOK, "Successfully signed in")
}
