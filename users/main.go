package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"users/db"
	"users/grpc"
	"users/handlers"
	"users/jwt"
	"users/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.ConnectToMongoDB(os.Getenv("MONGODB_URI")); err != nil {
		log.Fatal("Could not connect to MongoDB", err.Error())
	}
	if err := grpc.InitTasksClient(os.Getenv("TASKS_URI")); err != nil {
		log.Fatal("Could not init tasks service grpc client", err.Error())
	}

	privateFileFlag := flag.String("private", "", "path to JWT private key `file`")
	publicFileFlag := flag.String("public", "", "path to JWT public key `file`")
	portFlag := flag.Int("port", 8091, "http server port")
	flag.Parse()

	if err := jwt.LoadKeys(privateFileFlag, publicFileFlag); err != nil {
		log.Fatal("Could not load JWT keys: ", err.Error())
	}

	router := gin.Default()

	authGroup := router.Group("/auth")
	authGroup.POST("/signup", handlers.SignUp)
	authGroup.POST("/signin", handlers.SignIn)

	taskGroup := router.Group("/task").Use(middleware.JWTAuthMiddleware)
	taskGroup.POST("", handlers.CreateTask)
	taskGroup.GET("", handlers.ListTasks)
	taskGroup.GET("/:id", handlers.GetTask)
	taskGroup.PUT("/:id", handlers.UpdateTask)
	taskGroup.DELETE("/:id", handlers.DeleteTask)

	userGroup := router.Group("/task").Use(middleware.JWTAuthMiddleware)
	userGroup.PUT("user/:username", handlers.UpdateUser)

	log.Printf("Starting server on port %d", *portFlag)
	if err := router.Run(fmt.Sprintf(":%d", *portFlag)); err != nil {
		log.Fatal("Failed to run http server: ", err.Error())
	}
}
