package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"users/models"
)

var database *mongo.Database

func ConnectToMongoDB(uri string) error {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	database = client.Database("users")
	return err
}

func CreateUser(userCreds models.UserCredentials) error {
	userCreadsCollection := database.Collection("creds")
	_, err := userCreadsCollection.InsertOne(context.TODO(), userCreds)
	return err
}

func UpdateUserData(username string, userData models.UserData) error {
	userDataCollection := database.Collection("data")
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "data", Value: userData}}},
		{Key: "$setOnInsert", Value: filter},
	}
	opts := options.Update().SetUpsert(true)
	_, err := userDataCollection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func GetUserCredentials(username string) (models.UserCredentials, error) {
	userCreadsCollection := database.Collection("creds")
	filter := bson.D{{Key: "username", Value: username}}
	var userCreds models.UserCredentials
	err := userCreadsCollection.FindOne(context.TODO(), filter).Decode(&userCreds)
	if err != nil {
		return models.UserCredentials{}, err
	}
	return userCreds, nil
}
