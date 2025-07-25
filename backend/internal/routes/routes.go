package routes

import (
	"eduqr-backend/internal/controllers"
	"eduqr-backend/internal/middlewares"
	"time"

	"eduqr-backend/internal/database"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userController    *controllers.UserController
	eventController   *controllers.EventController
	roomController    *controllers.RoomController
	subjectController *controllers.SubjectController
	courseController  *controllers.CourseController
	authMiddleware    *middlewares.AuthMiddleware
}

func NewRouter(
	userController *controllers.UserController,
	eventController *controllers.EventController,
	roomController *controllers.RoomController,
	subjectController *controllers.SubjectController,
	courseController *controllers.CourseController,
	authMiddleware *middlewares.AuthMiddleware,
) *Router {
	return &Router{
		userController:    userController,
		eventController:   eventController,
		roomController:    roomController,
		subjectController: subjectController,
		courseController:  courseController,
		authMiddleware:    authMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
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
			// Profile routes (users can manage their own profile)
			users.GET("/profile", r.userController.GetProfile)
			users.PUT("/profile", r.userController.UpdateProfile)
			users.PUT("/profile/password", r.userController.ChangePassword)
			users.POST("/profile/validate-password", r.userController.ValidatePassword)

			// User management routes with role-based permissions
			users.GET("/all", r.userController.GetAllUsers)    // All authenticated users can view based on their role
			users.POST("/create", r.userController.CreateUser) // Only users who can manage roles

			// Parameterized routes with role-based permissions
			users.GET("/:id", r.userController.GetUserByID)           // View permissions based on role
			users.PUT("/:id", r.userController.UpdateUser)            // Manage permissions based on role
			users.PATCH("/:id/role", r.userController.UpdateUserRole) // Manage permissions based on role
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

		// Room routes (admin authentication required)
		rooms := v1.Group("/admin/rooms")
		rooms.Use(r.authMiddleware.AuthMiddleware())
		rooms.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			rooms.GET("", r.roomController.GetAllRooms)
			rooms.GET("/modular", r.roomController.GetModularRooms)
			rooms.POST("", r.roomController.CreateRoom)
			rooms.GET("/:id", r.roomController.GetRoomByID)
			rooms.PUT("/:id", r.roomController.UpdateRoom)
		}

		// Subject routes (admin authentication required)
		subjects := v1.Group("/admin/subjects")
		subjects.Use(r.authMiddleware.AuthMiddleware())
		subjects.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			subjects.GET("", r.subjectController.GetAllSubjects)
			subjects.POST("", r.subjectController.CreateSubject)
			subjects.GET("/:id", r.subjectController.GetSubjectByID)
			subjects.PUT("/:id", r.subjectController.UpdateSubject)
		}

		// Course routes (admin authentication required)
		courses := v1.Group("/admin/courses")
		courses.Use(r.authMiddleware.AuthMiddleware())
		courses.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			courses.GET("", r.courseController.GetAllCourses)
			courses.POST("", r.courseController.CreateCourse)
			courses.GET("/:id", r.courseController.GetCourseByID)
			courses.PUT("/:id", r.courseController.UpdateCourse)
			courses.GET("/by-date-range", r.courseController.GetCoursesByDateRange)
			courses.GET("/by-room/:roomId", r.courseController.GetCoursesByRoom)
			courses.GET("/by-teacher/:teacherId", r.courseController.GetCoursesByTeacher)
			courses.POST("/check-conflicts", r.courseController.CheckConflicts)
			courses.POST("/:id/check-conflicts", r.courseController.CheckConflictsForUpdate)
		}

		// Routes de suppression sécurisées
		deletionController := controllers.NewDeletionController(services.NewDeletionService(
			repositories.NewUserRepository(),
			repositories.NewCourseRepository(database.GetDB()),
			repositories.NewRoomRepository(database.GetDB()),
			repositories.NewSubjectRepository(),
		))

		// Suppression d'utilisateurs (Admin/Super Admin seulement)
		v1.DELETE("/admin/users/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("user"), deletionController.DeleteUser)

		// Suppression de salles (Admin/Super Admin seulement)
		v1.DELETE("/admin/rooms/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("room"), deletionController.DeleteRoom)

		// Suppression de matières (Admin/Super Admin seulement)
		v1.DELETE("/admin/subjects/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("subject"), deletionController.DeleteSubject)

		// Suppression de cours (Admin/Super Admin seulement)
		v1.DELETE("/admin/courses/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("course"), deletionController.DeleteCourse)
	}

	return router
}
