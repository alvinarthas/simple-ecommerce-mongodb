package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/models"
	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var err error

// RedirectHandler to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_G"),
			"clientSecret": os.Getenv("CLIENT_SECRET_G"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler Handle callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "error",
			"errors":  err.Error(),
		})
		c.Abort()
		return
	}

	var newUser = getOrRegisterUser(provider, user)
	var jtwToken = createToken(&newUser)

	c.JSON(200, gin.H{
		"success": true,
		"data":    newUser,
		"token":   jtwToken,
		"message": "berhasil login",
	})
}

// RegisterUser to store the new customer data into DB
func RegisterUser(c *gin.Context) {
	// Check Password confirmation
	password := c.PostForm("password")
	confirmedPassword := c.PostForm("confirmed_password")
	token, _ := RandomToken()

	// Return Error if not confirmed
	if password != confirmedPassword {
		c.JSON(500, gin.H{
			"success": false,
			"message": "password not confirmed",
			"errors":  err.Error()})
		c.Abort()
		return
	}

	// Hash the password
	hash, _ := HashPassword(password)

	// Get Form
	newUser := models.User{
		ID:                primitive.NewObjectID(),
		UserName:          c.PostForm("user_name"),
		FullName:          c.PostForm("full_name"),
		Email:             c.PostForm("email"),
		Password:          hash,
		VerificationToken: token,
	}

	collection := config.DB.Collection("users")

	_, err := collection.InsertOne(config.CTX, newUser)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "successfuly register user, please check your email",
		"data":   newUser,
	})
}

// LoginUser to get the token for access the system
func LoginUser(c *gin.Context) {
	// Get Login Form
	email := c.PostForm("email")
	password := c.PostForm("password")

	var userData models.User
	collection := config.DB.Collection("users")

	filter := bson.M{
		"email":     email,
		"is_active": 1,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&userData)

	if err == nil && CheckPasswordHash(password, userData.Password) {
		token := createToken(&userData)

		c.JSON(200, gin.H{
			"status": "success",
			"data":   userData,
			"token":  token,
		})
	} else {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "your email/password may be wrong",
		})
	}
}

// getOrRegisterUser new social ID to database
func getOrRegisterUser(provider string, user *structs.User) models.User {
	var userData models.User
	collection := config.DB.Collection("users")
	// Find One user

	filter := bson.M{
		"provider":  provider,
		"social_id": user.ID,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&userData)

	if err != nil {
		token, _ := RandomToken()

		newUser := models.User{
			ID:                primitive.NewObjectID(),
			FullName:          user.FullName,
			UserName:          user.Username,
			Email:             user.Email,
			SocialID:          user.ID,
			Provider:          provider,
			Avatar:            user.Avatar,
			VerificationToken: token,
		}

		_, err := collection.InsertOne(config.CTX, newUser)

		if err != nil {
			log.Fatal(err)
		}

		return newUser
	}

	return userData
}

// CreateToken to generate token for accesing the system
func createToken(user *models.User) string {
	// to send time expire, issue at (iat)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":        user.ID,
		"user_role":      user.Role,
		"user_store":     user.HaveStore,
		"store_username": user.Store.UserName,
		"exp":            time.Now().AddDate(0, 0, 1).Unix(), // expired after one day
		"iat":            time.Now().Unix(),                  // date created
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

// VerifyUserAccount to verify user and store account
func VerifyUserAccount(c *gin.Context) {
	var userData models.User
	collection := config.DB.Collection("users")

	verificationToken := c.Param("token")

	err = collection.FindOne(config.CTX, bson.M{"verification_token": verificationToken}).Decode(&userData)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "verification token mismatch"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": userData.ID}
	updateStatement := bson.M{"$set": bson.M{"is_active": 1}}

	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "verifiaction user error"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    result.ModifiedCount,
		"message": "User Successfully Verify",
	})
}

// VerifyStoreAccount to verify user and store account
func VerifyStoreAccount(c *gin.Context) {

	var userData models.User
	collection := config.DB.Collection("users")

	verificationToken := c.Param("token")

	err = collection.FindOne(config.CTX, bson.M{"store.verification_token": verificationToken}).Decode(&userData)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "verification token mismatch"})
		c.Abort()
		return
	}

	selector := bson.M{"_id": userData.ID}
	updateStatement := bson.M{"$set": bson.M{
		"is_active":       1,
		"have_store":      1,
		"store.is_active": 1,
	}}

	result, err := collection.UpdateOne(
		config.CTX,
		selector,
		updateStatement,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "verifiaction user error"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    result.ModifiedCount,
		"message": "User Successfully Verify",
	})
}
