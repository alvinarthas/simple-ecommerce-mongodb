package routes

import (
	"fmt"

	"github.com/alvinarthas/simple-ecommerce-mongodb/collections"
	"gopkg.in/mgo.v2/bson"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/gin-gonic/gin"
)

// TestFunc only for testing
func TestFunc(c *gin.Context) {
	var err error
	// INSERT

	// users := collections.User{
	// 	UserName: "alvinarthas",
	// 	FullName: "Alvin Arthas",
	// 	Email:    "arthas@alvin.com",
	// 	Password: "password",
	// 	SocialID: "28726294",
	// 	Provider: "github",
	// 	Avatar:   "https://avatars3.githubusercontent.com/u/28726294?v=4",
	// }

	collection := config.DB.Collection("users")

	// // // INSERT
	// insertResult, err := collection.InsertOne(config.CTX, users)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// FIND

	filter := bson.M{
		"social_id": "28726294",
	}

	var result collections.User

	err = collection.FindOne(config.CTX, filter).Decode(&result)
	if err != nil {
		fmt.Println("Error calling FindOne():", err)
	} else {
		fmt.Println("FindOne() result:", result)
	}
}
