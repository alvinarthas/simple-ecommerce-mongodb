package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/alvinarthas/simple-ecommerce-mongodb/collections"
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/danilopolani/gocialite/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
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
	// // Retrieve query params for state and code
	// state := c.Query("state")
	// code := c.Query("code")
	// provider := c.Param("provider")

	// // Handle callback and check for errors
	// user, _, err := config.Gocial.Handle(state, code)
	// if err != nil {
	// 	c.Writer.Write([]byte("Error: " + err.Error()))
	// 	return
	// }

	// // var newUser = getOrRegisterUser(provider, user)
	// // var jtwToken = createToken(&newUser)
}

// RegisterUser to store the new customer data into DB
func RegisterUser(c *gin.Context) {
}

// LoginUser to get the token for access the system
func LoginUser(c *gin.Context) {
}

// getOrRegisterUser new social ID to database
func getOrRegisterUser(provider string, user *structs.User) collections.User {
	var userData collections.User
	collection := config.DB.Collection("users")
	// Find One user

	filter := bson.M{
		"provider":  provider,
		"social_id": user.ID,
	}

	err = collection.FindOne(config.CTX, filter).Decode(&userData)

	if err == nil {
		token, _ := RandomToken()

		newUser := collections.User{
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
// func createToken(user *collections.User) string {
// var store collections.Store
// var storeID primitive.ObjectID

// if user.HaveStore == 1 {
// 	err = collection.FindOne(config.CTX, bson.M{}).Decode(&userData)
// 	// if config.DB.First(&store, "user_id = ?", user.ID).RecordNotFound() {
// 	// 	storeID = 0
// 	// }
// 	// storeID = store.ID
// }
// // to send time expire, issue at (iat)
// jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 	"user_id":    user.ID,
// 	"user_role":  user.Role,
// 	"user_store": user.HaveStore,
// 	"store_id":   storeID,
// 	"exp":        time.Now().AddDate(0, 0, 7).Unix(),
// 	"iat":        time.Now().Unix(),
// })

// // Sign and get the complete encoded token as a string using the secret
// tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

// if err != nil {
// 	fmt.Println(err)
// }

// return tokenString
// }
