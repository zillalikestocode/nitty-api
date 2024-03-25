package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty" validator:"required"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty" validator:"required"`
	Password string             `json:"password,omitempty" bson:"password,omitempty" validator:"required"`
}
