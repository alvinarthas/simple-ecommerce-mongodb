package main

import (
	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/alvinarthas/simple-ecommerce-mongodb/middleware"
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
		// Users
		users := apiV1.Group("/users")
		{
			// Normal Register and Login
			users.POST("/register", routes.RegisterUser)
			users.POST("/login", routes.LoginUser)

			// Social Auth or OAuth
			users.GET("/auth/:provider", routes.RedirectHandler)
			users.GET("/auth/:provider/callback", routes.CallbackHandler)

			// Verification Users
			users.GET("/verify/:token", routes.VerifyUserAccount)
		}
		// Store
		stores := apiV1.Group("/stores")
		{
			// Registration
			stores.POST("/register", middleware.IsAuth(), routes.RegisterStore)
			//show all store products & Account Info
			stores.GET("/:username", routes.GetStore)
			stores.GET("/:username/info", middleware.HaveStore(), routes.InfoStore)
			// Verification User and Store Account
			stores.GET("/verify/:token", routes.VerifyStoreAccount)
		}

		// Product CRUD by Store
		products := apiV1.Group("/products")
		{
			products.GET("/", routes.GetAllProducts) // every product in every store
			products.GET("/:slug", routes.GetProduct)
			products.POST("/", middleware.HaveStore(), routes.CreateProduct)
			products.PUT("/:slug", middleware.HaveStore(), routes.UpdateProduct)
			products.DELETE("/:slug", middleware.HaveStore(), routes.DeleteProduct)
		}

		// Category
		categories := apiV1.Group("/categories")
		{
			// Initilize Http method for Category Crud
			categories.GET("/", routes.GetAllCategories)
			categories.GET("/:slug", middleware.IsAdmin(), routes.GetCategory)
			categories.GET("/:slug/products", routes.GetCategoryProduct)
			categories.POST("/", middleware.IsAdmin(), routes.CreateCategory)
			categories.PUT("/:slug", middleware.IsAdmin(), routes.UpdateCategory)
			categories.DELETE("/:slug", middleware.IsAdmin(), routes.DeleteCategory)
		}

		apiV1.GET("/testfunc", routes.TestFunc)
	}

	router.Run()
}
