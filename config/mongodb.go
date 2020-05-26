package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var CTX = context.TODO()

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

	DB = client.Database("test")
}
