package routes

import (
	"fmt"
	"log"

	"github.com/alvinarthas/simple-ecommerce-mongodb/collections"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/gin-gonic/gin"
)

// TestFunc only for testing
func TestFunc(c *gin.Context) {
	var err error
	// INSERT

	store := collections.Store{
		ID:       primitive.NewObjectID(),
		Name:     "Toko Bagus",
		UserName: "bagus_toko",
		Adress:   "Cipa Cipa",
		Email:    "tokok@bagus.com",
		Phone:    "password",
		Avatar:   "https://avatars3.githubusercontent.com/u/28726294?v=4",
	}

	users := collections.User{
		ID:       primitive.NewObjectID(),
		UserName: "cipa_)ipa",
		FullName: "Cipa Cipa",
		Email:    "cipa@alvin.com",
		Password: "password",
		SocialID: "2872629224",
		Provider: "github",
		Avatar:   "https://avatars3.githubusercontent.com/u/28726294?v=4",
		Store:    store,
	}

	collection := config.DB.Collection("users")

	// // // INSERT
	insertResult, err := collection.InsertOne(config.CTX, users)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// FIND ONE

	// filter := bson.M{
	// 	"social_id": "28726294",
	// }

	// var result collections.User
	// err = collection.FindOne(config.CTX, filter).Decode(&result)

	// if err != nil {
	// 	fmt.Println("Error calling FindOne():", err)
	// } else {
	// 	fmt.Println("FindOne() result:", result)
	// }

	// // Pass these options to the Find method
	// findOptions := options.Find()

	// // Here's an array in which you can store the decoded documents
	// var results []*collections.User

	// // Passing bson.D{{}} as the filter matches all documents in the collection
	// cur, err := collection.Find(config.CTX, filter, findOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Finding multiple documents returns a cursor
	// // Iterating through the cursor allows us to decode documents one at a time
	// for cur.Next(context.TODO()) {

	// 	// create a value into which the single document can be decoded
	// 	var elem collections.User
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	results = append(results, &elem)
	// }

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Close the cursor once finished
	// cur.Close(config.CTX)

	c.JSON(200, gin.H{
		"status": "successfuly register user, please check your email",
		"data":   users,
	})
}
