package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB to set the global variable
var DB *mongo.Database

// CTX to set the context, to know the time of using
var CTX = context.TODO()

// InitDB to initialize the MongoDB connection
func InitDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(CTX, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(CTX, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("simple_ecommerce")
}
