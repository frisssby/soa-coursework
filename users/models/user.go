package models

import "time"

type UserCredentials struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type UserData struct {
	FirstName string    `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string    `bson:"last_name,omitempty" json:"last_name,omitempty"`
	DateBirth time.Time `bson:"date_birth,omitempty" json:"date_birth,omitempty"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Phone     string    `bson:"phone,omitempty" json:"phone,omitempty"`
}
