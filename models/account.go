package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Account collection
type Account struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Avatar      string             `bson:"avatar" json:"avatar"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description" json:"description"`
	Account     string             `bson:"account" json:"account"`
}
