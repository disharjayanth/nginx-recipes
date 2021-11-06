package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string
}
