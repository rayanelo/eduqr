package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"fmt"
	"log"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = database.AutoMigrate(&models.User{}, &models.AuditLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository()

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)

	// Create users
	log.Println("Creating users...")

	// Super Admin
	superAdminReq := &models.CreateUserRequest{
		Email:           "superadmin@eduqr.com",
		Password:        "superadmin123",
		ConfirmPassword: "superadmin123",
		FirstName:       "Super",
		LastName:        "Admin",
		Phone:           "+1234567890",
		Address:         "123 Admin Street",
		Role:            "super_admin",
	}

	if _, err := userService.CreateUser(superAdminReq); err != nil {
		log.Println("Failed to create Super Admin:", err)
	} else {
		log.Println("Super Admin created successfully")
	}

	// Admin
	adminReq := &models.CreateUserRequest{
		Email:           "admin@eduqr.com",
		Password:        "admin123",
		ConfirmPassword: "admin123",
		FirstName:       "Admin",
		LastName:        "User",
		Phone:           "+1234567891",
		Address:         "456 Admin Avenue",
		Role:            "admin",
	}

	if _, err := userService.CreateUser(adminReq); err != nil {
		log.Println("Failed to create Admin:", err)
	} else {
		log.Println("Admin created successfully")
	}

	// Professeur 1
	prof1Req := &models.CreateUserRequest{
		Email:           "prof1@eduqr.com",
		Password:        "prof123",
		ConfirmPassword: "prof123",
		FirstName:       "Jean",
		LastName:        "Dupont",
		Phone:           "+1234567892",
		Address:         "789 Teacher Street",
		Role:            "professeur",
	}

	if _, err := userService.CreateUser(prof1Req); err != nil {
		log.Println("Failed to create Professeur prof1@eduqr.com:", err)
	} else {
		log.Println("Professeur prof1@eduqr.com created successfully")
	}

	// Professeur 2
	prof2Req := &models.CreateUserRequest{
		Email:           "prof2@eduqr.com",
		Password:        "prof123",
		ConfirmPassword: "prof123",
		FirstName:       "Marie",
		LastName:        "Martin",
		Phone:           "+1234567893",
		Address:         "321 Teacher Avenue",
		Role:            "professeur",
	}

	if _, err := userService.CreateUser(prof2Req); err != nil {
		log.Println("Failed to create Professeur prof2@eduqr.com:", err)
	} else {
		log.Println("Professeur prof2@eduqr.com created successfully")
	}

	// Etudiant 1
	etudiant1Req := &models.CreateUserRequest{
		Email:           "etudiant1@eduqr.com",
		Password:        "student123",
		ConfirmPassword: "student123",
		FirstName:       "Pierre",
		LastName:        "Durand",
		Phone:           "+1234567894",
		Address:         "654 Student Street",
		Role:            "etudiant",
	}

	if _, err := userService.CreateUser(etudiant1Req); err != nil {
		log.Println("Failed to create Etudiant etudiant1@eduqr.com:", err)
	} else {
		log.Println("Etudiant etudiant1@eduqr.com created successfully")
	}

	// Etudiant 2
	etudiant2Req := &models.CreateUserRequest{
		Email:           "etudiant2@eduqr.com",
		Password:        "student123",
		ConfirmPassword: "student123",
		FirstName:       "Sophie",
		LastName:        "Leroy",
		Phone:           "+1234567895",
		Address:         "987 Student Avenue",
		Role:            "etudiant",
	}

	if _, err := userService.CreateUser(etudiant2Req); err != nil {
		log.Println("Failed to create Etudiant etudiant2@eduqr.com:", err)
	} else {
		log.Println("Etudiant etudiant2@eduqr.com created successfully")
	}

	// Etudiant 3
	etudiant3Req := &models.CreateUserRequest{
		Email:           "etudiant3@eduqr.com",
		Password:        "student123",
		ConfirmPassword: "student123",
		FirstName:       "Lucas",
		LastName:        "Moreau",
		Phone:           "+1234567896",
		Address:         "147 Student Boulevard",
		Role:            "etudiant",
	}

	if _, err := userService.CreateUser(etudiant3Req); err != nil {
		log.Println("Failed to create Etudiant etudiant3@eduqr.com:", err)
	} else {
		log.Println("Etudiant etudiant3@eduqr.com created successfully")
	}

	// Créer des logs d'audit de test
	if err := createTestAuditLogs(); err != nil {
		log.Printf("Erreur lors de la création des logs d'audit de test: %v", err)
	} else {
		log.Println("✅ Logs d'audit de test créés avec succès")
	}

	log.Println("🎉 Seeding terminé avec succès!")

	log.Println("\nTest accounts:")
	log.Println("Super Admin: superadmin@eduqr.com / superadmin123")
	log.Println("Admin: admin@eduqr.com / admin123")
	log.Println("Professeur: prof1@eduqr.com / prof123")
	log.Println("Etudiant: etudiant1@eduqr.com / student123")
}

// createTestAuditLogs crée des logs d'audit de test
func createTestAuditLogs() error {
	db := database.GetDB()

	// Récupérer quelques utilisateurs pour créer des logs
	var users []models.User
	if err := db.Limit(3).Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return fmt.Errorf("aucun utilisateur trouvé pour créer des logs d'audit")
	}

	// Actions de test
	testActions := []struct {
		action       string
		resourceType string
		description  string
	}{
		{models.ActionLogin, "", "Connexion utilisateur"},
		{models.ActionCreate, models.ResourceRoom, "Création d'une nouvelle salle"},
		{models.ActionUpdate, models.ResourceUser, "Modification d'un utilisateur"},
		{models.ActionCreate, models.ResourceSubject, "Création d'une nouvelle matière"},
		{models.ActionUpdate, models.ResourceCourse, "Modification d'un cours"},
		{models.ActionDelete, models.ResourceEvent, "Suppression d'un événement"},
		{models.ActionCreate, models.ResourceAbsence, "Création d'une absence"},
		{models.ActionUpdate, models.ResourceAbsence, "Validation d'une absence"},
	}

	// Créer des logs pour chaque utilisateur
	for i, user := range users {
		for j, testAction := range testActions {
			// Créer des dates variées (derniers 7 jours)
			createdAt := time.Now().Add(-time.Duration(i*24+j*3) * time.Hour)

			auditLog := models.AuditLog{
				UserID:       user.ID,
				UserEmail:    user.Email,
				UserRole:     user.Role,
				Action:       testAction.action,
				ResourceType: testAction.resourceType,
				ResourceID:   nil, // Pas d'ID spécifique pour les tests
				Description:  testAction.description,
				OldValues:    nil,
				NewValues:    nil,
				IPAddress:    "127.0.0.1",
				UserAgent:    "Test-Agent/1.0",
				CreatedAt:    createdAt,
			}

			if err := db.Create(&auditLog).Error; err != nil {
				return fmt.Errorf("erreur lors de la création du log d'audit: %v", err)
			}
		}
	}

	return nil
}
