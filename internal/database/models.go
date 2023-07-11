package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	PID         primitive.ObjectID `json:"pid" bson:"_id"`
	Header      string             `json:"header" bson:"header"`
	Description string             `json:"description" bson:"description"`
	Author      string             `json:"author" bson:"author"`
}
