package handlers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllCouriers is to get all category -> Admin Only
func GetAllCouriers(c *gin.Context) {
	collection := config.DB.Collection("couriers")

	findOptions := options.Find()
	var results []*models.Courier

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
		var elem models.Courier
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

// GetCourier is to show courier detail
func GetCourier(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("couriers")
	var courier models.Courier

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&courier)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   courier,
	})
}

// CreateCourier is store courier data
func CreateCourier(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("couriers")
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

	costRaw, _ := strconv.ParseInt(c.PostForm("cost"), 10, 64)
	cost := int32(costRaw)
	// Get Store Data
	item := models.Courier{
		ID:          primitive.NewObjectID(),
		Name:        c.PostForm("name"),
		Slug:        slug,
		Description: c.PostForm("description"),
		Avatar:      c.PostForm("avatar"),
		Cost:        cost,
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
		"message": "Store Courier Successful",
		"data":    item,
	})
}

// UpdateCourier is to update courier
func UpdateCourier(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("courier")
	var courier models.Courier

	// Get Parameter
	getSlug := c.Param("slug")

	filter := bson.M{
		"slug": getSlug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&courier)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	if c.PostForm("name") != courier.Name {
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

	costRaw, _ := strconv.ParseInt(c.PostForm("cost"), 10, 64)
	cost := int32(costRaw)

	selector := bson.M{"_id": courier.ID}
	updateStatement := bson.M{"$set": bson.M{
		"name":        c.PostForm("name"),
		"description": c.PostForm("description"),
		"slug":        getSlug,
		"avatar":      c.PostForm("avatar"),
		"cost":        cost,
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

// DeleteCourier is to delete courier
func DeleteCourier(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("courier")
	var courier models.Courier

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err := collection.FindOne(config.CTX, filter).Decode(&courier)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": courier.ID}

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
		"message": "Delete Courier Success",
	})
}
