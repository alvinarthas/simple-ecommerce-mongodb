package main

import (
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	//  Setting Default Router
	router := gin.Default()

	// Initialize Version
	apiV1 := router.Group("/api/v1/")
	{
		// Normal Register and Login
		apiV1.GET("/register", routes.RegisterUser)
		apiV1.POST("/login", routes.LoginUser)
	}

	router.Run()
}
