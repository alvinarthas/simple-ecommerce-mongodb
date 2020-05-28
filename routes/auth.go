package routes

import (
	"net/http"
	"os"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/gin-gonic/gin"
)

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

// // Register new social ID to database
// func getOrRegisterUser(provider string, user *structs.User) collections.User {
// 	var userData collections.User
// }
