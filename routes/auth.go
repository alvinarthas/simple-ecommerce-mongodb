package routes

import (
	"fmt"
	"log"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/model"
	"github.com/gin-gonic/gin"
)

// RegisterUser to store the new customer data into DB
func RegisterUser(c *gin.Context) {
	ash := model.Trainer{"Milea", 10, "Pallet Town"}

	collection := config.DB.Collection("trainers")

	// INSERT
	insertResult, err := collection.InsertOne(config.CTX, ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

// LoginUser to get the token for access the system
func LoginUser(c *gin.Context) {
}
