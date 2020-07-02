package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Address collection
type Address struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Phone    string             `bson:"phone" json:"phone"`
	City     string             `bson:"city" json:"city"`
	PostCode string             `bson:"post_code" json:"post_code"`
	Address  string             `bson:"address" json:"address"`
}
