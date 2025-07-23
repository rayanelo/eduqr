package routes

import (
	"eduqr-backend/internal/controllers"
	"eduqr-backend/internal/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userController  *controllers.UserController
	eventController *controllers.EventController
	authMiddleware  *middlewares.AuthMiddleware
}

func NewRouter(
	userController *controllers.UserController,
	eventController *controllers.EventController,
	authMiddleware *middlewares.AuthMiddleware,
) *Router {
	return &Router{
		userController:  userController,
		eventController: eventController,
		authMiddleware:  authMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "EduQR API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userController.Register)
			auth.POST("/login", r.userController.Login)
		}

		// User routes (authentication required)
		users := v1.Group("/users")
		users.Use(r.authMiddleware.AuthMiddleware())
		{
			users.GET("/profile", r.userController.GetProfile)
			users.PUT("/profile", r.userController.UpdateProfile)
			users.GET("/:id", r.userController.GetUserByID)
		}

		// Event routes (authentication required)
		events := v1.Group("/events")
		events.Use(r.authMiddleware.AuthMiddleware())
		{
			events.GET("", r.eventController.GetUserEvents)
			events.POST("", r.eventController.CreateEvent)
			events.GET("/range", r.eventController.GetEventsByDateRange)
			events.GET("/:id", r.eventController.GetEventByID)
			events.PUT("/:id", r.eventController.UpdateEvent)
			events.DELETE("/:id", r.eventController.DeleteEvent)
		}
	}

	return router
}
