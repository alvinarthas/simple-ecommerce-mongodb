package main

import (
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/handlers"
	"github.com/alvinarthas/simple-ecommerce-mongodb/middleware"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	// Start the MongoDB Connection and Load the Environtment Variable
	config.InitDB()
	gotenv.Load()

	//  Setting Default Router
	router := gin.Default()

	// Initialize Version
	apiV1 := router.Group("/api/v1/")
	{
		// Users
		users := apiV1.Group("/users")
		{
			// Normal Register and Login
			users.POST("/register", handlers.RegisterUser)
			users.POST("/login", handlers.LoginUser)

			// Social Auth or OAuth
			users.GET("/auth/:provider", handlers.RedirectHandler)
			users.GET("/auth/:provider/callback", handlers.CallbackHandler)

			// Verification Users
			users.PATCH("/verify/:token", handlers.VerifyUserAccount)
		}
		// Store
		stores := apiV1.Group("/stores")
		{
			//show all store products & Account Info
			stores.GET("/:username", handlers.GetStore)
			stores.GET("/:username/info", middleware.HaveStore(), handlers.InfoStore)

			// Registration
			stores.POST("/register", middleware.IsAuth(), handlers.RegisterStore)

			// Verification User and Store Account
			stores.PATCH("/verify/:token", handlers.VerifyStoreAccount)
		}

		// Product CRUD
		products := apiV1.Group("/products")
		{
			products.GET("/", handlers.GetAllProducts) // every product in every store
			products.GET("/:slug", handlers.GetProduct)
			products.POST("/", middleware.HaveStore(), handlers.CreateProduct)
			products.PUT("/:slug", middleware.HaveStore(), handlers.UpdateProduct)
			products.DELETE("/:slug", middleware.HaveStore(), handlers.DeleteProduct)
		}

		// Category CRUD
		categories := apiV1.Group("/categories")
		{
			// Initilize Http method for Category Crud
			categories.GET("/", handlers.GetAllCategories)
			categories.GET("/:slug", middleware.IsAdmin(), handlers.GetCategory)
			categories.GET("/:slug/products", handlers.GetCategoryProduct)
			categories.POST("/", middleware.IsAdmin(), handlers.CreateCategory)
			categories.PUT("/:slug", middleware.IsAdmin(), handlers.UpdateCategory)
			categories.DELETE("/:slug", middleware.IsAdmin(), handlers.DeleteCategory)
		}

		apiV1.GET("/testfunc", handlers.TestFunc)
	}

	// Run The Server
	router.Run()
}
