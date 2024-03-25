package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	Admin bool               `json:"admin" bson:"admin"`
}

type Announcement struct {
	ID      primitive.ObjectID `json:"id" bson:"id"`
	Creator struct {
		Name string             `json:"name" bson:"name"`
		ID   primitive.ObjectID `json:"id" bson:"id"`
	}
	Date    primitive.DateTime `json:"date" bson:"date"`
	Message string             `json:"message" bson:"message"`
}

type Event struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Date        primitive.DateTime `json:"date" bson:"date"`
	Time        string             `json:"time" bson:"time"`
	Address     string             `json:"address" bson:"address"`
}

type Community struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty" validator:"required"`
	Description   string             `json:"description,omitempty" bson:"description,omitempty" validator:"required"`
	Owner         primitive.ObjectID `json:"owner,omitempty" bson:"owner,omitempty" validator:"required"`
	Members       []Member           `json:"members,omitempty" bson:"members,omitempty"`
	Announcements []Announcement     `json:"announcements,omitempty" bson:"announcements,omitempty"`
	Events        []Event            `json:"events,omitempty" bson:"events,omitempty"`
}
