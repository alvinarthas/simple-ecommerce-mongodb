package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
Product is Belong To Store
One Store can have many products
One Product can only have one Store
One Product can only have one Categories
*/

// Product collection
type Product struct {
	ID           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Slug         string             `bson:"slug" json:"slug"`
	Description  string             `bson:"description" json:"description"`
	Condition    int                `bson:"condition" json:"condition"`
	Price        int                `bson:"price" json:"price"`
	InitialStock int                `bson:"initial_stock" json:"initial_stock"`
	Weight       int                `bson:"weight" json:"weight"`
	Category     string             `bson:"category" json:"category"`
	Store        string             `bson:"store" json:"store"`
}
