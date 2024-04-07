package models

import "time"

type UserCredentials struct {
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
}

type UserData struct {
	FirstName string    `bson:"first_name" json:"first_name" binding:"required"`
	LastName  string    `bson:"last_name" json:"last_name" binding:"required"`
	DateBirth time.Time `bson:"date_birth" json:"date_birth" binding:"required"`
	Email     string    `bson:"email" json:"email" binding:"required"`
	Phone     string    `bson:"phone" json:"phone" binding:"required"`
}
