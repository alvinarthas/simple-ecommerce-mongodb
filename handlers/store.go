package handlers

import (
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	User can have more than one store
	User still can shop and User still can sell their things in their stores
*/

// RegisterStore is to register
func RegisterStore(c *gin.Context) {

	// Get User ID from Authorization token & Initialization
	UserobjID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))
	token, _ := RandomToken()
	collection := config.DB.Collection("users")

	store := models.Store{
		Name:              c.PostForm("name"),
		UserName:          c.PostForm("user_name"),
		Adress:            c.PostForm("adress"),
		Email:             c.PostForm("email"),
		Phone:             c.PostForm("phone"),
		Avatar:            c.PostForm("avatar"),
		VerificationToken: token,
	}

	selector := bson.M{"_id": UserobjID}
	updateStatement := bson.M{"$set": bson.M{
		"store":      store,
		"have_store": 1,
	}}

	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "registration error"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    result.ModifiedCount,
		"message": "User Store Registration Success",
	})
}

// GetStore to gett all products that the store have
func GetStore(c *gin.Context) {
	// // Get Store ID from Authorization token
	// storeUsername := c.Param("username")

	// // Set Query Params for Filtering & Sorting
	// queryCategory := c.Query("category")
	// querySort := c.Query("sort")
	// queryPriceMin := c.Query("pricemin")
	// queryPriceMax := c.Query("pricemax")
	// queryCondition := c.Query("condition")

	// // Initialization
	// collection := config.DB.Collection("users")

	// // Dynamic query
	// filter := bson.M{}

}

// InfoStore to get store account info
func InfoStore(c *gin.Context) {
	// Get Store ID from Authorization token
	storeUsername := c.Param("username")
	collection := config.DB.Collection("users")

	filter := bson.M{
		"store.user_name": storeUsername,
	}

	var users models.User
	err = collection.FindOne(config.CTX, filter).Decode(&users)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   users.Store,
	})
}
