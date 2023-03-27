package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notes struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	User  primitive.ObjectID `bson:"user"`
	Title string             `json:"title" bson:"title"`
	Data  []string           `json:"data" bson:"data"`
	Tag   string             `json:"tag" bson:"tag"`
	Date  string             `json:"date" bson:"date"`
}
