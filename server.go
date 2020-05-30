package main

import (
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	config.InitDB()
	gotenv.Load()

	//  Setting Default Router
	router := gin.Default()

	// Initialize Version
	apiV1 := router.Group("/api/v1/")
	{
		// Social Auth or OAuth
		apiV1.GET("/auth/:provider", routes.RedirectHandler)
		apiV1.GET("/auth/:provider/callback", routes.CallbackHandler)

		// Normal Register and Login
		apiV1.GET("/register", routes.RegisterUser)
		apiV1.POST("/login", routes.LoginUser)

		// Verification User and Store Account
		apiV1.GET("/verify/store/:token", routes.VerifyStoreAccount)
		apiV1.GET("/verify/user/:token", routes.VerifyUserAccount)

		apiV1.GET("/testfunc", routes.TestFunc)
	}

	router.Run()
}
