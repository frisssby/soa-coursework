package db

import (
	"context"
	"log"

	"tasks/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func ConnectToMongoDB(uri string) error {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	database = client.Database("tasks")
	return err
}

func CreateTask(task *models.Task) (*models.Task, error) {
	taskCollection := database.Collection("tasks")
	res, err := taskCollection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID).MarshalText()
	newTask := &models.Task{
		TaskId:      string(id),
		UserId:      task.UserId,
		Description: task.Description,
	}
	task.TaskId = string(id)
	return newTask, nil
}

func UpdateTask(id string, task models.Task) error {
	taskCollection := database.Collection("tasks")
	var objectId primitive.ObjectID
	if err := objectId.UnmarshalText([]byte(id)); err != nil {
		return mongo.ErrNoDocuments
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{
		{Key: "$set", Value: task},
		{Key: "$setOnInsert", Value: filter},
	}
	opts := options.Update().SetUpsert(true)
	_, err := taskCollection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func DeleteTask(id string) error {
	taskCollection := database.Collection("tasks")
	var objectId primitive.ObjectID
	if err := objectId.UnmarshalText([]byte(id)); err != nil {
		return mongo.ErrNoDocuments
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err := taskCollection.DeleteOne(context.Background(), filter)
	return err
}

func GetTask(id string) (*models.Task, error) {
	taskCollection := database.Collection("tasks")
	var objectId primitive.ObjectID
	if err := objectId.UnmarshalText([]byte(id)); err != nil {
		return nil, mongo.ErrNoDocuments
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	var task models.Task
	err := taskCollection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func ListsTasks(userId string, pageSize, pageId int) ([]models.Task, error) {
	taskCollection := database.Collection("tasks")
	l := int64(pageSize)
	skip := int64(pageSize * pageId)
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}
	filter := bson.D{{Key: "user_id", Value: userId}}
	curr, err := taskCollection.Find(context.Background(), filter, &fOpt)
	if err != nil {
		return nil, err
	}
	result := make([]models.Task, 0)
	for curr.Next(context.Background()) {
		var el models.Task
		if err := curr.Decode(&el); err != nil {
			log.Println(err)
		} else {
			result = append(result, el)
		}
	}
	return result, nil
}
