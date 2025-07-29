package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eduqr-backend/config"
	"eduqr-backend/internal/controllers"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/middlewares"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/routes"
	"eduqr-backend/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// Auto migrate models
	if err := database.AutoMigrate(&models.User{}, &models.Event{}, &models.Room{}, &models.Subject{}, &models.Course{}, &models.AuditLog{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	eventRepo := repositories.NewEventRepository()
	roomRepo := repositories.NewRoomRepository(database.GetDB())
	subjectRepo := repositories.NewSubjectRepository()
	courseRepo := repositories.NewCourseRepository(database.GetDB())
	auditLogRepo := repositories.NewAuditLogRepository()

	// Parse JWT expiration
	jwtExpiration, err := time.ParseDuration(cfg.JWT.Expiration)
	if err != nil {
		log.Fatalf("Failed to parse JWT expiration: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, jwtExpiration)
	eventService := services.NewEventService(eventRepo)
	roomService := services.NewRoomService(roomRepo)
	subjectService := services.NewSubjectService(subjectRepo)
	courseService := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)
	auditLogService := services.NewAuditLogService(auditLogRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	eventController := controllers.NewEventController(eventService)
	roomController := controllers.NewRoomController(roomService)
	subjectController := controllers.NewSubjectController(subjectService)
	courseController := controllers.NewCourseController(courseService)
	auditLogController := controllers.NewAuditLogController(auditLogService)

	// Initialize middleware
	authMiddleware := middlewares.NewAuthMiddleware(cfg.JWT.Secret)
	auditMiddleware := middlewares.NewAuditMiddleware(auditLogService)

	// Initialize router
	router := routes.NewRouter(userController, eventController, roomController, subjectController, courseController, auditLogController, authMiddleware, auditMiddleware)
	app := router.SetupRoutes()

	// Create server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: app,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
