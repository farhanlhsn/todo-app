package main

import (
	"todo-app/controllers"
	"todo-app/helpers"
	"todo-app/initializers"
	"todo-app/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	println("Application ready to use")
}

func main() {
	router := gin.Default()

	// Add global middlewares
	router.Use(middlewares.CORS())
	router.Use(middlewares.GeneralRateLimit())

	// Health check endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, helpers.FormatResponseWithoutData(true, "Welcome to the Todo App API v1.0!"))
	})

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Health check for v1
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, helpers.FormatResponseWithoutData(true, "Todo App API v1.0 is healthy!"))
		})

		// Authentication routes with stricter rate limiting
		auth := v1.Group("/auth")
		auth.Use(middlewares.AuthRateLimit()) // Apply strict rate limiting for auth
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.LogIn)
			auth.POST("/logout", middlewares.RequiredAuth, controllers.LogOut)
		}

		// Task routes
		task := v1.Group("/task")
		{
			task.Use(middlewares.RequiredAuth)

			// Basic CRUD
			task.POST("/", controllers.CreateTask)
			task.GET("/", controllers.GetTasks)
			task.GET("/:id", controllers.GetTask)
			task.PUT("/:id", controllers.UpdateTask)
			task.DELETE("/:id", controllers.DeleteTask)

			// Task Status Operations
			task.PUT("/:id/complete", controllers.CompletedTask)
			task.PUT("/:id/uncomplete", controllers.UncompletedTask)

			// Task Filtering
			task.GET("/completed", controllers.GetCompletedTasks)
			task.GET("/pending", controllers.GetPendingTasks)
			task.GET("/overdue", controllers.GetOverdueTasks)

			// Task Search & Filter
			task.GET("/search", controllers.SearchTasks)

			// Category Management
			taskCategories := task.Group("/categories")
			{
				taskCategories.GET("/", controllers.GetCategories)
				taskCategories.POST("/", controllers.CreateCategory)
				taskCategories.GET("/:name", controllers.GetTasksByCategory)
				taskCategories.DELETE("/:id", controllers.DeleteCategory)
			}
		}

		// User routes
		user := v1.Group("/users")
		{
			user.Use(middlewares.RequiredAuth)
			user.GET("/profile", controllers.Profile)
			user.GET("/stats", controllers.Stats)
		}
	}

	router.Run()
}
