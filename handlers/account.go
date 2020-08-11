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

// GetAllAccounts is to get all accounts -> Admin Only
func GetAllAccounts(c *gin.Context) {
	collection := config.DB.Collection("accounts")

	findOptions := options.Find()
	var results []*models.Account

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
		var elem models.Account
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	// Return JSON
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Berhasil menampilkan semua data Bank Account",
		"data":    results,
	})
}

// GetAccount is to show account detail
func GetAccount(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("accounts")
	var account models.Account

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&account)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   account,
	})
}

// CreateAccount is store account data
func CreateAccount(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("accounts")
	// Get Parameter
	slug := slug.Make(c.PostForm("name") + "-" + c.PostForm("account"))

	// Check Slug
	filterCheck := bson.M{
		"slug": slug,
	}

	_, err := collection.Find(config.CTX, filterCheck)

	if err != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// Get Store Data
	item := models.Account{
		ID:          primitive.NewObjectID(),
		Name:        c.PostForm("name"),
		Slug:        slug,
		Description: c.PostForm("description"),
		Avatar:      c.PostForm("avatar"),
		Account:     c.PostForm("account"),
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
		"message": "Store Account Successful",
		"data":    item,
	})
}

// UpdateAccount is to update account
func UpdateAccount(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("accounts")
	var account models.Account

	// Get Parameter
	getSlug := c.Param("slug")

	filter := bson.M{
		"slug": getSlug,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&account)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	if c.PostForm("name") != account.Name || c.PostForm("account") != account.Account {
		// Make New Slug
		newSlug := slug.Make(c.PostForm("name") + "-" + c.PostForm("account"))
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

	selector := bson.M{"_id": account.ID}
	updateStatement := bson.M{"$set": bson.M{
		"name":        c.PostForm("name"),
		"description": c.PostForm("description"),
		"slug":        getSlug,
		"avatar":      c.PostForm("avatar"),
		"account":     c.PostForm("account"),
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

// DeleteAccount is to delete account
func DeleteAccount(c *gin.Context) {
	// Initialization
	collection := config.DB.Collection("accounts")
	var account models.Account

	// Get Parameter
	slug := c.Param("slug")

	filter := bson.M{
		"slug": slug,
	}

	err := collection.FindOne(config.CTX, filter).Decode(&account)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": account.ID}

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
		"message": "Delete Account Success",
	})
}
