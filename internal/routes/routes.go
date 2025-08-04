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
	userController     *controllers.UserController
	eventController    *controllers.EventController
	roomController     *controllers.RoomController
	subjectController  *controllers.SubjectController
	courseController   *controllers.CourseController
	auditLogController *controllers.AuditLogController
	absenceController  *controllers.AbsenceController
	presenceController *controllers.PresenceController
	authMiddleware     *middlewares.AuthMiddleware
	auditMiddleware    *middlewares.AuditMiddleware
}

func NewRouter(
	userController *controllers.UserController,
	eventController *controllers.EventController,
	roomController *controllers.RoomController,
	subjectController *controllers.SubjectController,
	courseController *controllers.CourseController,
	auditLogController *controllers.AuditLogController,
	absenceController *controllers.AbsenceController,
	presenceController *controllers.PresenceController,
	authMiddleware *middlewares.AuthMiddleware,
	auditMiddleware *middlewares.AuditMiddleware,
) *Router {
	return &Router{
		userController:     userController,
		eventController:    eventController,
		roomController:     roomController,
		subjectController:  subjectController,
		courseController:   courseController,
		auditLogController: auditLogController,
		absenceController:  absenceController,
		presenceController: presenceController,
		authMiddleware:     authMiddleware,
		auditMiddleware:    auditMiddleware,
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
			auth.POST("/login", r.auditMiddleware.AuditLoginMiddleware(), r.userController.Login)
		}

		// User routes (authentication required)
		users := v1.Group("/users")
		users.Use(r.authMiddleware.AuthMiddleware())
		{
			// Profile routes (users can manage their own profile)
			users.GET("/profile", r.userController.GetProfile)
			users.PUT("/profile", r.auditMiddleware.AuditMiddleware("update", "user"), r.userController.UpdateProfile)
			users.PUT("/profile/password", r.userController.ChangePassword)
			users.POST("/profile/validate-password", r.userController.ValidatePassword)

			// User management routes with role-based permissions
			users.GET("/all", r.userController.GetAllUsers)                                                         // All authenticated users can view based on their role
			users.POST("/create", r.auditMiddleware.AuditMiddleware("create", "user"), r.userController.CreateUser) // Only users who can manage roles

			// Parameterized routes with role-based permissions
			users.GET("/:id", r.userController.GetUserByID)                                                                // View permissions based on role
			users.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "user"), r.userController.UpdateUser)            // Manage permissions based on role
			users.PATCH("/:id/role", r.auditMiddleware.AuditMiddleware("update", "user"), r.userController.UpdateUserRole) // Manage permissions based on role
		}

		// Event routes (authentication required)
		events := v1.Group("/events")
		events.Use(r.authMiddleware.AuthMiddleware())
		{
			events.GET("", r.eventController.GetUserEvents)
			events.POST("", r.auditMiddleware.AuditMiddleware("create", "event"), r.eventController.CreateEvent)
			events.GET("/range", r.eventController.GetEventsByDateRange)
			events.GET("/:id", r.eventController.GetEventByID)
			events.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "event"), r.eventController.UpdateEvent)
			events.DELETE("/:id", r.auditMiddleware.AuditMiddleware("delete", "event"), r.eventController.DeleteEvent)
		}

		// Absence routes (authentication required)
		absences := v1.Group("/absences")
		absences.Use(r.authMiddleware.AuthMiddleware())
		{
			// Routes pour tous les utilisateurs authentifiés
			absences.POST("", r.auditMiddleware.AuditMiddleware("create", "absence"), r.absenceController.CreateAbsence)            // Étudiants seulement
			absences.GET("/my", r.absenceController.GetMyAbsences)                                                                  // Étudiants seulement
			absences.GET("/teacher", r.absenceController.GetTeacherAbsences)                                                        // Professeurs seulement
			absences.GET("/stats", r.absenceController.GetAbsenceStats)                                                             // Tous selon leur rôle
			absences.GET("/filter", r.absenceController.GetAbsencesWithFilters)                                                     // Admins et professeurs
			absences.GET("/:id", r.absenceController.GetAbsenceByID)                                                                // Selon les permissions
			absences.POST("/:id/review", r.auditMiddleware.AuditMiddleware("update", "absence"), r.absenceController.ReviewAbsence) // Professeurs et admins
			absences.DELETE("/:id", r.auditMiddleware.AuditMiddleware("delete", "absence"), r.absenceController.DeleteAbsence)      // Selon les permissions
		}

		// Presence routes (authentication required)
		presences := v1.Group("/presences")
		presences.Use(r.authMiddleware.AuthMiddleware())
		{
			// Routes pour les étudiants
			presences.POST("/scan", r.auditMiddleware.AuditMiddleware("create", "presence"), r.presenceController.ScanQRCode) // Étudiants seulement
			presences.GET("/my", r.presenceController.GetMyPresences)                                                         // Étudiants seulement

			// Routes pour les professeurs et admins
			presences.GET("/course/:courseId", r.presenceController.GetPresencesByCourse)                                                                              // Professeurs et admins
			presences.GET("/course/:courseId/stats", r.presenceController.GetPresenceStats)                                                                            // Professeurs et admins
			presences.POST("/course/:courseId/create-all", r.auditMiddleware.AuditMiddleware("create", "presence"), r.presenceController.CreatePresenceForAllStudents) // Professeurs et admins
		}

		// QR Code routes (authentication required)
		qrCodes := v1.Group("/qr-codes")
		qrCodes.Use(r.authMiddleware.AuthMiddleware())
		{
			// Routes pour les professeurs et admins
			qrCodes.GET("/course/:courseId", r.presenceController.GetQRCodeInfo)                                                                        // Professeurs et admins
			qrCodes.POST("/course/:courseId/regenerate", r.auditMiddleware.AuditMiddleware("update", "qr_code"), r.presenceController.RegenerateQRCode) // Professeurs et admins
		}

		// Room routes (admin authentication required)
		rooms := v1.Group("/admin/rooms")
		rooms.Use(r.authMiddleware.AuthMiddleware())
		rooms.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			rooms.GET("", r.roomController.GetAllRooms)
			rooms.GET("/modular", r.roomController.GetModularRooms)
			rooms.POST("", r.auditMiddleware.AuditMiddleware("create", "room"), r.roomController.CreateRoom)
			rooms.GET("/:id", r.roomController.GetRoomByID)
			rooms.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "room"), r.roomController.UpdateRoom)
		}

		// Subject routes (admin authentication required)
		subjects := v1.Group("/admin/subjects")
		subjects.Use(r.authMiddleware.AuthMiddleware())
		subjects.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			subjects.GET("", r.subjectController.GetAllSubjects)
			subjects.POST("", r.auditMiddleware.AuditMiddleware("create", "subject"), r.subjectController.CreateSubject)
			subjects.GET("/:id", r.subjectController.GetSubjectByID)
			subjects.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "subject"), r.subjectController.UpdateSubject)
		}

		// Course routes (admin authentication required)
		courses := v1.Group("/admin/courses")
		courses.Use(r.authMiddleware.AuthMiddleware())
		courses.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			courses.GET("", r.courseController.GetAllCourses)
			courses.POST("", r.auditMiddleware.AuditMiddleware("create", "course"), r.courseController.CreateCourse)
			courses.GET("/:id", r.courseController.GetCourseByID)
			courses.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "course"), r.courseController.UpdateCourse)
			courses.GET("/by-date-range", r.courseController.GetCoursesByDateRange)
			courses.GET("/by-room/:roomId", r.courseController.GetCoursesByRoom)
			courses.GET("/by-teacher/:teacherId", r.courseController.GetCoursesByTeacher)
			courses.POST("/check-conflicts", r.courseController.CheckConflicts)
			courses.POST("/:id/check-conflicts", r.courseController.CheckConflictsForUpdate)
		}

		// Public course routes (authentication required, no admin role required)
		publicCourses := v1.Group("/courses")
		publicCourses.Use(r.authMiddleware.AuthMiddleware())
		{
			publicCourses.GET("", r.courseController.GetAllCourses) // Lecture seule pour tous les utilisateurs
			publicCourses.GET("/:id", r.courseController.GetCourseByID)
			publicCourses.GET("/by-teacher/:teacherId", r.courseController.GetCoursesByTeacher)
			publicCourses.GET("/by-room/:roomId", r.courseController.GetCoursesByRoom)

			// Routes pour les professeurs (création et modification de leurs propres cours)
			publicCourses.POST("", r.auditMiddleware.AuditMiddleware("create", "course"), r.courseController.CreateCourse)
			publicCourses.PUT("/:id", r.auditMiddleware.AuditMiddleware("update", "course"), r.courseController.UpdateCourse)
			publicCourses.DELETE("/:id", r.auditMiddleware.AuditMiddleware("delete", "course"), r.courseController.DeleteCourse)
		}

		// Admin absence routes (admin authentication required)
		adminAbsences := v1.Group("/admin/absences")
		adminAbsences.Use(r.authMiddleware.AuthMiddleware())
		adminAbsences.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			adminAbsences.GET("", r.absenceController.GetAllAbsences)
		}

		// Admin presence routes (admin authentication required)
		adminPresences := v1.Group("/admin/presences")
		adminPresences.Use(r.authMiddleware.AuthMiddleware())
		adminPresences.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			adminPresences.GET("", r.presenceController.GetPresencesWithFilters)
		}

		// Routes de suppression sécurisées
		deletionController := controllers.NewDeletionController(services.NewDeletionService(
			repositories.NewUserRepository(),
			repositories.NewCourseRepository(database.GetDB()),
			repositories.NewRoomRepository(database.GetDB()),
			repositories.NewSubjectRepository(),
		))

		// Suppression d'utilisateurs (Admin/Super Admin seulement)
		v1.DELETE("/admin/users/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("user"), r.auditMiddleware.AuditMiddleware("delete", "user"), deletionController.DeleteUser)

		// Suppression de salles (Admin/Super Admin seulement)
		v1.DELETE("/admin/rooms/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("room"), r.auditMiddleware.AuditMiddleware("delete", "room"), deletionController.DeleteRoom)

		// Suppression de matières (Admin/Super Admin seulement)
		v1.DELETE("/admin/subjects/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("subject"), r.auditMiddleware.AuditMiddleware("delete", "subject"), deletionController.DeleteSubject)

		// Suppression de cours (Admin/Super Admin seulement)
		v1.DELETE("/admin/courses/:id", r.authMiddleware.AuthMiddleware(), r.authMiddleware.CanDeleteMiddleware("course"), r.auditMiddleware.AuditMiddleware("delete", "course"), deletionController.DeleteCourse)

		// Audit Log routes (Admin/Super Admin only)
		auditLogs := v1.Group("/admin/audit-logs")
		auditLogs.Use(r.authMiddleware.AuthMiddleware())
		auditLogs.Use(r.authMiddleware.RoleMiddleware("admin"))
		{
			auditLogs.GET("", r.auditLogController.GetAuditLogs)
			auditLogs.GET("/stats", r.auditLogController.GetAuditLogStats)
			auditLogs.GET("/recent", r.auditLogController.GetRecentAuditLogs)
			auditLogs.GET("/:id", r.auditLogController.GetAuditLogByID)
			auditLogs.GET("/user/:user_id/activity", r.auditLogController.GetUserActivity)
			auditLogs.GET("/resource/:resource_type/:resource_id", r.auditLogController.GetResourceHistory)
			auditLogs.DELETE("/clean", r.auditLogController.CleanOldLogs)
		}
	}

	return router
}
