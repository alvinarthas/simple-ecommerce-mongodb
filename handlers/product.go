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

// GetAllProducts to
func GetAllProducts(c *gin.Context) {
	collection := config.DB.Collection("products")

	findOptions := options.Find()
	var results []*models.Product

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
		var elem models.Product
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
	var product models.Product

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
	item := models.Product{
		ID:           primitive.NewObjectID(),
		Name:         c.PostForm("name"),
		Slug:         slug,
		Description:  c.PostForm("description"),
		Price:        price,
		Condition:    condition,
		InitialStock: stock,
		Weight:       weight,
		Category:     c.PostForm("category"),
		Store:        c.MustGet("jwt_user_id").(string),
	}

	insertResult, _ := collection.InsertOne(config.CTX, item)

	// Save product object id to Store
	userCollection := config.DB.Collection("users")

	userID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))
	selector := bson.M{"_id": userID}
	updateStatement := bson.M{"$push": bson.M{"store.products": insertResult.InsertedID}}

	_, err = userCollection.UpdateOne(
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
		"status":  "success",
		"message": "Store Product Successful",
		"data":    item,
	})

}

// UpdateProduct to
func UpdateProduct(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	var product models.Product

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
	userCollection := config.DB.Collection("users")

	userID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))
	filterUser := bson.M{
		"_id":            userID,
		"store.products": product.ID,
	}

	productCheck, _ := userCollection.CountDocuments(config.CTX, filterUser)

	if productCheck == 1 {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "This Data is Forbidden"})
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
	var product models.Product

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

	// Delete Product from user collection
	userCollection := config.DB.Collection("users")

	userID, _ := primitive.ObjectIDFromHex(c.MustGet("jwt_user_id").(string))
	userSelector := bson.M{"_id": userID}
	updateStatement := bson.M{"$pull": bson.M{"store.products": product.ID}}

	_, err = userCollection.UpdateOne(
		config.CTX,
		userSelector,
		updateStatement,
	)

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

// ProductSearch with parameter can do dinamically based on the including query param
func ProductSearch(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("products")
	var results []*models.Product
	filter := bson.M{}
	findOptions := options.Find()

	// Get Parameter
	keyword := c.Query("keyword")
	minPrice, _ := strconv.Atoi(c.Query("minPrice"))
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
	condition, _ := strconv.Atoi(c.Query("condition"))
	category := c.Query("category")
	// sort := c.Query("sort")

	// IF There is keyword
	if len(keyword) > 0 {
		filter["$or"] = []interface{}{
			bson.M{"name": primitive.Regex{ // ^ is start with, $ is end with
				Pattern: keyword,
				Options: "i",
			}},
			bson.M{"store": primitive.Regex{
				Pattern: keyword,
				Options: "i",
			}},
		}
	}

	// Include Price Filters
	if len(c.Query("minPrice")) > 0 && len(c.Query("maxPrice")) > 0 {
		filter["price"] = bson.M{
			"$gte": minPrice,
			"$lte": maxPrice,
		}
	} else if len(c.Query("minPrice")) > 0 {
		filter["price"] = bson.M{
			"$lte": maxPrice,
		}
	} else if len(c.Query("maxPrice")) > 0 {
		filter["price"] = bson.M{
			"$gte": minPrice,
		}
	}

	// Include Condition Filter
	if len(c.Query("condition")) > 0 {
		filter["condition"] = condition
	}

	// Include Category Filter
	if len(category) > 0 {
		filter["category"] = category
	}

	// findOptions.SetSort(bson.D{{"name", -1}})

	cur, err := collection.Find(config.CTX, filter, findOptions)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	for cur.Next(config.CTX) {
		// create a value into which the single document can be decoded
		var elem models.Product
		err := cur.Decode(&elem)
		if err != nil {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "record not found"})
			c.Abort()
			return
		}

		results = append(results, &elem)
	}

	// Return JSON
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Show Products Data",
		"data":    results,
	})
}
