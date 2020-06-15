package routes

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/alvinarthas/simple-ecommerce-mongodb/collections"
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// GetAllProducts to
func GetAllProducts(c *gin.Context) {
	collection := config.DB.Collection("products")

	findOptions := options.Find()
	var results []*collections.Product

	cur, err := collection.Find(config.CTX, bson.M{}, findOptions)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	for cur.Next(config.CTX) {
		// create a value into which the single document can be decoded
		var elem collections.Product
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	// Return JSON
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Berhasil menampilkan semua data kategori",
		"data":    results,
	})
}

// GetProduct to
func GetProduct(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	var product collections.Product

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&product)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   product,
	})
}

// CreateProduct to
func CreateProduct(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	// Get Parameter
	slug := slug.Make(c.PostForm("name"))

	// Check Slug
	filterCheck := bson.M{
		"slug": slug,
	}

	_, err := collection.Find(config.CTX, filterCheck)

	if err != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// set Parameter POST
	price, _ := strconv.Atoi(c.PostForm("price"))
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	weight, _ := strconv.Atoi(c.PostForm("weight"))
	condition, _ := strconv.Atoi(c.PostForm("condition"))

	// Get Store Data
	item := collections.Product{
		ID:           primitive.NewObjectID(),
		Name:         c.PostForm("name"),
		Slug:         slug,
		Description:  c.PostForm("description"),
		Price:        price,
		Condition:    condition,
		InitialStock: stock,
		Weight:       weight,
		Category:     c.PostForm("category"),
	}

	_, err = collection.InsertOne(config.CTX, item)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Store Product Successful",
		"data":    item,
	})

}

// UpdateProduct to
func UpdateProduct(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	var product collections.Product

	// Get Parameter
	getSlug := c.Param("slug")

	filter := bson.M{
		"slug": getSlug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&product)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	// To make sure, it is the right store account who do the update
	storeID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_store_id").(string))
	if storeID != product.ID {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "this data is forbidden"})
		c.Abort()
		return
	}

	// Get Form for updating the product
	price, _ := strconv.Atoi(c.PostForm("price"))
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	weight, _ := strconv.Atoi(c.PostForm("weight"))
	condition, _ := strconv.Atoi(c.PostForm("condition"))

	if c.PostForm("name") != product.Name {
		// Make New Slug
		newSlug := slug.Make(c.PostForm("name"))
		// Check Slug
		filterCheck := bson.M{
			"slug": newSlug,
		}

		slugCount, _ := collection.CountDocuments(config.CTX, filterCheck)

		if slugCount > 0 {
			getSlug = getSlug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
		} else {
			getSlug = newSlug
		}
	}

	selector := bson.M{"_id": product.ID}
	updateStatement := bson.M{"$set": bson.M{
		"name":         c.PostForm("name"),
		"description":  c.PostForm("description"),
		"slug":         getSlug,
		"price":        price,
		"condition":    stock,
		"intial_stock": weight,
		"weight":       condition,
		"category":     c.PostForm("category"),
	}}

	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   result,
	})
}

// DeleteProduct to
func DeleteProduct(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	var product collections.Product

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err := collection.FindOne(config.CTX, filter).Decode(&product)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": product.ID}

	_, err = collection.DeleteOne(context.TODO(), selector)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Delete Product Success",
	})
}