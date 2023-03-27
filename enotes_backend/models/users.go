package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" binding:"required,min=3,max=50"`
	Email    string             `json:"email" bson:"email" binding:"required,email,min=5,max=100" unique:"true"`
	Password string             `json:"password" bson:"password" binding:"required"`
	Date     string             `json:"date,omitempty" bson:"date"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}