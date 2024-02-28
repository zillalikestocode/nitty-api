package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Community struct {
	ID            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty" validator:"required"`
	Description   string               `json:"description,omitempty" bson:"description,omitempty" validator:"required"`
	Owner         primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty" validator:"required"`
	Members       []primitive.ObjectID `json:"members,omitempty" bson:"members,omitempty"`
	Announcements []struct {
		ID      primitive.ObjectID `json:"id" bson:"id"`
		Creator struct {
			Name string             `json:"name" bson:"name"`
			ID   primitive.ObjectID `json:"id" bson:"id"`
		}
		Date    primitive.DateTime `json:"date" bson:"date"`
		Message string             `json:"message" bson:"message"`
	} `json:"announcements,omitempty" bson:"announcements,omitempty"`
	Events []struct {
		ID          primitive.ObjectID `json:"id" bson:"id"`
		Name        string             `json:"name" bson:"name"`
		Description string             `json:"description" bson:"description"`
		Date        primitive.DateTime `json:"date" bson:"date"`
		Time        string             `json:"time" bson:"time"`
	} `json:"events,omitempty" bson:"events,omitempty"`
}
