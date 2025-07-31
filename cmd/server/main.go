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
	"eduqr-backend/pkg/utils"

	"gorm.io/gorm"
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
	if err := database.AutoMigrate(&models.User{}, &models.Event{}, &models.Room{}, &models.Subject{}, &models.Course{}, &models.AuditLog{}, &models.Absence{}, &models.Presence{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Créer le super admin par défaut
	createDefaultSuperAdmin()

	// Créer les données de base
	createDefaultData()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	eventRepo := repositories.NewEventRepository()
	roomRepo := repositories.NewRoomRepository(database.GetDB())
	subjectRepo := repositories.NewSubjectRepository()
	courseRepo := repositories.NewCourseRepository(database.GetDB())
	auditLogRepo := repositories.NewAuditLogRepository()
	absenceRepo := repositories.NewAbsenceRepository(database.GetDB())
	presenceRepo := repositories.NewPresenceRepository(database.GetDB())

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
	absenceService := services.NewAbsenceService(absenceRepo, courseRepo, userRepo)
	presenceService := services.NewPresenceService(presenceRepo, courseRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	eventController := controllers.NewEventController(eventService)
	roomController := controllers.NewRoomController(roomService)
	subjectController := controllers.NewSubjectController(subjectService)
	courseController := controllers.NewCourseController(courseService)
	auditLogController := controllers.NewAuditLogController(auditLogService)
	absenceController := controllers.NewAbsenceController(absenceService)
	presenceController := controllers.NewPresenceController(presenceService)

	// Initialize middleware
	authMiddleware := middlewares.NewAuthMiddleware(cfg.JWT.Secret)
	auditMiddleware := middlewares.NewAuditMiddleware(auditLogService)

	// Initialize router
	router := routes.NewRouter(userController, eventController, roomController, subjectController, courseController, auditLogController, absenceController, presenceController, authMiddleware, auditMiddleware)
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

// createDefaultSuperAdmin crée un super admin par défaut s'il n'existe pas
func createDefaultSuperAdmin() {
	db := database.GetDB()

	// Vérifier si un super admin existe déjà
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)

	if count > 0 {
		log.Println("Super admin existe déjà, pas de création")
		return
	}

	// Créer le super admin par défaut
	superAdmin := models.User{
		Email:        "admin@eduqr.com",
		ContactEmail: "admin@eduqr.com",
		Password:     func() string { p, _ := utils.HashPassword("Admin123!"); return p }(),
		FirstName:    "Super",
		LastName:     "Administrateur",
		Phone:        "+33123456789",
		Address:      "123 Rue de l'Administration, 75001 Paris",
		Role:         models.RoleSuperAdmin,
	}

	if err := db.Create(&superAdmin).Error; err != nil {
		log.Printf("Erreur lors de la création du super admin: %v", err)
		return
	}

	log.Println("✅ Super admin créé avec succès!")
	log.Println("   Email: admin@eduqr.com")
	log.Println("   Mot de passe: Admin123!")
}

// createDefaultData crée les données de base (salles, matières)
func createDefaultData() {
	db := database.GetDB()

	// Créer les salles par défaut
	createDefaultRooms(db)

	// Créer les matières par défaut
	createDefaultSubjects(db)
}

// createDefaultRooms crée les salles par défaut
func createDefaultRooms(db *gorm.DB) {
	var count int64
	db.Model(&models.Room{}).Count(&count)

	if count > 0 {
		log.Println("Salles existent déjà, pas de création")
		return
	}

	// Créer la salle modulable parent
	modularRoom := models.Room{
		Name:      "Salle Modulable A",
		Building:  "Bâtiment Principal",
		Floor:     "1er étage",
		IsModular: true,
	}

	if err := db.Create(&modularRoom).Error; err != nil {
		log.Printf("Erreur lors de la création de la salle modulable: %v", err)
		return
	}

	// Créer les sous-salles
	subRooms := []models.Room{
		{
			Name:      "Salle Modulable A1",
			Building:  "Bâtiment Principal",
			Floor:     "1er étage",
			IsModular: false,
			ParentID:  &modularRoom.ID,
		},
		{
			Name:      "Salle Modulable A2",
			Building:  "Bâtiment Principal",
			Floor:     "1er étage",
			IsModular: false,
			ParentID:  &modularRoom.ID,
		},
	}

	for _, room := range subRooms {
		if err := db.Create(&room).Error; err != nil {
			log.Printf("Erreur lors de la création de la sous-salle %s: %v", room.Name, err)
		}
	}

	// Créer la salle normale
	normalRoom := models.Room{
		Name:      "Salle Normale B",
		Building:  "Bâtiment Principal",
		Floor:     "2ème étage",
		IsModular: false,
	}

	if err := db.Create(&normalRoom).Error; err != nil {
		log.Printf("Erreur lors de la création de la salle normale: %v", err)
		return
	}

	log.Println("✅ Salles créées avec succès!")
}

// createDefaultSubjects crée les matières par défaut
func createDefaultSubjects(db *gorm.DB) {
	var count int64
	db.Model(&models.Subject{}).Count(&count)

	if count > 0 {
		log.Println("Matières existent déjà, pas de création")
		return
	}

	subjects := []models.Subject{
		{
			Name:        "Mathématiques",
			Code:        "MATH101",
			Description: "Cours de mathématiques fondamentales",
		},
		{
			Name:        "Informatique",
			Code:        "INFO101",
			Description: "Introduction à la programmation",
		},
		{
			Name:        "Physique",
			Code:        "PHYS101",
			Description: "Physique générale",
		},
		{
			Name:        "Anglais",
			Code:        "ANG101",
			Description: "Anglais technique",
		},
		{
			Name:        "Économie",
			Code:        "ECO101",
			Description: "Principes d'économie",
		},
	}

	for _, subject := range subjects {
		if err := db.Create(&subject).Error; err != nil {
			log.Printf("Erreur lors de la création de la matière %s: %v", subject.Name, err)
		}
	}

	log.Println("✅ Matières créées avec succès!")
}
