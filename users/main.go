package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"users/db"
	"users/jwt"
)

func main() {
	if err := db.ConnectToMongoDB(os.Getenv("MONGODB_URI")); err != nil {
		log.Fatal("Could not connect to MongoDB", err.Error())
	}

	privateFile := flag.String("private", "", "path to JWT private key `file`")
	publicFile := flag.String("public", "", "path to JWT public key `file`")
	port := flag.Int("port", 8091, "http server port")
	flag.Parse()

	if err := jwt.LoadKeys(privateFile, publicFile); err != nil {
		log.Fatal("Could not load JWT keys: ", err.Error())
	}

	router := gin.Default()
	router.POST("auth/signup", signUp)
	router.POST("auth/signin", signIn)
	router.PUT("user/:username", updateUser)
	router.Run(fmt.Sprintf(":%d", *port))
}
