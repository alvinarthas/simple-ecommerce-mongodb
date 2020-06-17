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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// GetAllCategories is to get all category -> Admin Only
func GetAllCategories(c *gin.Context) {
	collection := config.DB.Collection("categories")

	findOptions := options.Find()
	var results []*models.Category

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
		var elem models.Category
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

// GetCategoryProduct tp get products of the category
func GetCategoryProduct(c *gin.Context) {

	// Initialization
	collection := config.DB.Collection("products")
	var product models.Product

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"category": slug,
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

// GetCategory is to get spesific product -> Store
func GetCategory(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("categories")
	var category models.Category

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&category)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   category,
	})
}

// CreateCategory is to create new category
func CreateCategory(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("categories")
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

	// Get Store Data
	item := models.Category{
		ID:          primitive.NewObjectID(),
		Name:        c.PostForm("name"),
		Slug:        slug,
		Description: c.PostForm("description"),
		Icon:        c.PostForm("icon"),
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
		"message": "Store Category Successful",
		"data":    item,
	})
}

// UpdateCategory is to update existing product -> Store
func UpdateCategory(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("categories")
	var category models.Category

	// Get Parameter
	getSlug := c.Param("slug")

	filter := bson.M{
		"slug": getSlug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&category)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	if c.PostForm("name") != category.Name {
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

	selector := bson.M{"_id": category.ID}
	updateStatement := bson.M{"$set": bson.M{
		"name":        c.PostForm("name"),
		"description": c.PostForm("description"),
		"slug":        getSlug,
		"icon":        c.PostForm("icon"),
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

// DeleteCategory is to delete existing category
func DeleteCategory(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("categories")
	var category models.Category

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err := collection.FindOne(config.CTX, filter).Decode(&category)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": category.ID}

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
		"message": "Delete Category Success",
	})
}
