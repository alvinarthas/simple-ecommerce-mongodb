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
			// users.GET("/notifications", handlers.TestFunc)

			// Verification Users
			users.PATCH("/verify/:token", handlers.VerifyUserAccount)

			// User Carts
			carts := users.Group("/carts")
			{
				carts.GET("/", middleware.IsAuth(), handlers.GetAllCarts) // every product in every store
				carts.POST("/", middleware.IsAuth(), handlers.AddCarts)
				carts.PUT("/:slug", middleware.IsAuth(), handlers.EditCart)
				carts.DELETE("/:slug", middleware.IsAuth(), handlers.DeleteCart)
			}
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

		// Courier CRUD
		couriers := apiV1.Group("/couriers")
		{
			// Initilize Http method for Courier Crud
			couriers.GET("/", handlers.GetAllCouriers)
			couriers.GET("/:slug", middleware.IsAdmin(), handlers.GetCourier)
			couriers.POST("/", middleware.IsAdmin(), handlers.CreateCourier)
			couriers.PUT("/:slug", middleware.IsAdmin(), handlers.UpdateCourier)
			couriers.DELETE("/:slug", middleware.IsAdmin(), handlers.DeleteCourier)
		}

		// Courier Account
		accounts := apiV1.Group("/accounts")
		{
			// Initilize Http method for E-Commerce Account Crud
			accounts.GET("/", handlers.GetAllAccounts)
			accounts.GET("/:slug", middleware.IsAdmin(), handlers.GetAccount)
			accounts.POST("/", middleware.IsAdmin(), handlers.CreateAccount)
			accounts.PUT("/:slug", middleware.IsAdmin(), handlers.UpdateAccount)
			accounts.DELETE("/:slug", middleware.IsAdmin(), handlers.DeleteAccount)
		}

		apiV1.GET("/search", handlers.ProductSearch)

		// Personal Testing bug and experiment
		apiV1.GET("/testfunc", handlers.TestFunc)
	}

	// Run The Server
	router.Run()
}
