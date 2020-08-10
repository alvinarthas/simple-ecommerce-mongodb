package handlers

import (
	"log"
	"strconv"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllCarts is to get all saved carts
func GetAllCarts(c *gin.Context) {
	// Get User ID from Authorization token & Initialization
	UserobjID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))

	// Variable Initialization
	collection := config.DB.Collection("users")

	// Check Slug
	var results []*models.Cart
	selector := bson.M{"_id": UserobjID}
	findOptions := options.Find()

	cur, err := collection.Find(config.CTX, selector, findOptions)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	for cur.Next(config.CTX) {
		// create a value into which the single document can be decoded
		var elem models.Cart
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	// Return JSON
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Berhasil menampilkan semua data cart",
		"data":    results,
	})
}

// AddCarts is add new cart
func AddCarts(c *gin.Context) {
	// Get User ID from Authorization token & Initialization
	UserobjID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))

	// Set Variables
	collection := config.DB.Collection("users")
	qty, _ := strconv.Atoi(c.PostForm("qty"))

	selector := bson.M{
		"_id": UserobjID,
	}

	// Add into data stuct
	cart := models.Cart{
		Product:     c.PostForm("product"),
		Qty:         qty,
		Description: c.PostForm("description"),
	}

	updateStatement := bson.M{"$set": bson.M{
		"cart": cart,
	}}

	// Insert into collection Users
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

// EditCart is to update cart
func EditCart(c *gin.Context) {
	// Get User ID from Authorization token & Initialization
	UserobjID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))

	// Set Variables
	collection := config.DB.Collection("users")
	qty, _ := strconv.Atoi(c.PostForm("qty"))
	product := c.PostForm("product")

	selector := bson.M{
		"_id":          UserobjID,
		"cart.product": product,
	}

	updateStatement := bson.M{"$set": bson.M{
		"cart.product":     product,
		"cart.qty":         qty,
		"cart.description": c.PostForm("description"),
	}}

	// Insert into collection Users
	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "update error"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    result.ModifiedCount,
		"message": "User Cart Update Success",
	})
}

// DeleteCart is to delete cart
func DeleteCart(c *gin.Context) {
	// Get User ID from Authorization token & Initialization
	UserobjID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))

	// Set Variables
	collection := config.DB.Collection("users")
	product := c.PostForm("product")

	selector := bson.M{
		"_id":          UserobjID,
		"cart.product": product,
	}

	updateStatement := bson.M{"$pull": bson.M{
		"cart.product": product,
	}}

	// Insert into collection Users
	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "update error"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    result.ModifiedCount,
		"message": "User Cart Delete Success",
	})
}
