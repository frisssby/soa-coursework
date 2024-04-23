package models

type Task struct {
	TaskId      string `bson:"_id,omitempty"`
	UserId      string `bson:"user_id"`
	Description string `bson:"description"`
	Status      string `bson:"status"`
}
