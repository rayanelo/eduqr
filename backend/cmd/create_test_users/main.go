package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
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
	err = database.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)

	log.Println("🚀 Création des utilisateurs de test...")

	// Créer des utilisateurs de test avec des mots de passe simples
	testUsers := []models.CreateUserRequest{
		{
			Email:           "test_superadmin@eduqr.com",
			Password:        "test123",
			ConfirmPassword: "test123",
			FirstName:       "Test",
			LastName:        "SuperAdmin",
			Phone:           "+1234567890",
			Address:         "123 Test Street",
			Role:            "super_admin",
		},
		{
			Email:           "test_admin@eduqr.com",
			Password:        "test123",
			ConfirmPassword: "test123",
			FirstName:       "Test",
			LastName:        "Admin",
			Phone:           "+1234567891",
			Address:         "456 Test Avenue",
			Role:            "admin",
		},
		{
			Email:           "test_prof@eduqr.com",
			Password:        "test123",
			ConfirmPassword: "test123",
			FirstName:       "Test",
			LastName:        "Professeur",
			Phone:           "+1234567892",
			Address:         "789 Test Boulevard",
			Role:            "professeur",
		},
		{
			Email:           "test_student@eduqr.com",
			Password:        "test123",
			ConfirmPassword: "test123",
			FirstName:       "Test",
			LastName:        "Étudiant",
			Phone:           "+1234567893",
			Address:         "321 Test Road",
			Role:            "etudiant",
		},
		{
			Email:           "test_student2@eduqr.com",
			Password:        "test123",
			ConfirmPassword: "test123",
			FirstName:       "Test",
			LastName:        "Étudiant2",
			Phone:           "+1234567894",
			Address:         "654 Test Lane",
			Role:            "etudiant",
		},
	}

	createdCount := 0
	for _, userReq := range testUsers {
		_, err := userService.CreateUser(&userReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de l'utilisateur %s: %v", userReq.Email, err)
		} else {
			createdCount++
			log.Printf("✅ Utilisateur créé: %s (%s) - %s", userReq.Email, userReq.FirstName+" "+userReq.LastName, userReq.Role)
		}
	}

	log.Printf("🎉 %d utilisateurs de test créés avec succès!", createdCount)
	log.Printf("")
	log.Printf("📋 Comptes de test disponibles:")
	log.Printf("   - Super Admin: test_superadmin@eduqr.com / test123")
	log.Printf("   - Admin: test_admin@eduqr.com / test123")
	log.Printf("   - Professeur: test_prof@eduqr.com / test123")
	log.Printf("   - Étudiant: test_student@eduqr.com / test123")
	log.Printf("   - Étudiant 2: test_student2@eduqr.com / test123")
	log.Printf("")
	log.Printf("✅ Vous pouvez maintenant utiliser ces comptes dans vos tests!")
}
