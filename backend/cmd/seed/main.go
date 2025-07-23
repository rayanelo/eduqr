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

	// Initialize database
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

	// Create Super Admin user
	superAdmin := &models.CreateUserRequest{
		Email:           "superadmin@eduqr.com",
		Password:        "superadmin123",
		ConfirmPassword: "superadmin123",
		FirstName:       "Super",
		LastName:        "Admin",
		Phone:           "+1234567890",
		Address:         "123 Admin Street",
		Role:            models.RoleSuperAdmin,
	}

	_, err = userService.CreateUser(superAdmin)
	if err != nil {
		if err.Error() == "user already exists" {
			fmt.Println("Super Admin already exists")
		} else {
			log.Printf("Failed to create Super Admin: %v", err)
		}
	} else {
		fmt.Println("Super Admin created successfully")
	}

	// Create Admin user
	admin := &models.CreateUserRequest{
		Email:           "admin@eduqr.com",
		Password:        "admin123",
		ConfirmPassword: "admin123",
		FirstName:       "Admin",
		LastName:        "User",
		Phone:           "+1234567891",
		Address:         "456 Admin Avenue",
		Role:            models.RoleAdmin,
	}

	_, err = userService.CreateUser(admin)
	if err != nil {
		if err.Error() == "user already exists" {
			fmt.Println("Admin already exists")
		} else {
			log.Printf("Failed to create Admin: %v", err)
		}
	} else {
		fmt.Println("Admin created successfully")
	}

	// Create Professeur users
	professeurs := []*models.CreateUserRequest{
		{
			Email:           "prof1@eduqr.com",
			Password:        "prof123",
			ConfirmPassword: "prof123",
			FirstName:       "Jean",
			LastName:        "Dupont",
			Phone:           "+1234567892",
			Address:         "789 Teacher Street",
			Role:            models.RoleProfesseur,
		},
		{
			Email:           "prof2@eduqr.com",
			Password:        "prof123",
			ConfirmPassword: "prof123",
			FirstName:       "Marie",
			LastName:        "Martin",
			Phone:           "+1234567893",
			Address:         "321 Teacher Avenue",
			Role:            models.RoleProfesseur,
		},
	}

	for _, prof := range professeurs {
		_, err = userService.CreateUser(prof)
		if err != nil {
			if err.Error() == "user already exists" {
				fmt.Printf("Professeur %s already exists\n", prof.Email)
			} else {
				log.Printf("Failed to create Professeur %s: %v", prof.Email, err)
			}
		} else {
			fmt.Printf("Professeur %s created successfully\n", prof.Email)
		}
	}

	// Create Etudiant users
	etudiants := []*models.CreateUserRequest{
		{
			Email:           "etudiant1@eduqr.com",
			Password:        "student123",
			ConfirmPassword: "student123",
			FirstName:       "Pierre",
			LastName:        "Durand",
			Phone:           "+1234567894",
			Address:         "654 Student Street",
			Role:            models.RoleEtudiant,
		},
		{
			Email:           "etudiant2@eduqr.com",
			Password:        "student123",
			ConfirmPassword: "student123",
			FirstName:       "Sophie",
			LastName:        "Leroy",
			Phone:           "+1234567895",
			Address:         "987 Student Avenue",
			Role:            models.RoleEtudiant,
		},
		{
			Email:           "etudiant3@eduqr.com",
			Password:        "student123",
			ConfirmPassword: "student123",
			FirstName:       "Lucas",
			LastName:        "Moreau",
			Phone:           "+1234567896",
			Address:         "147 Student Road",
			Role:            models.RoleEtudiant,
		},
	}

	for _, etudiant := range etudiants {
		_, err = userService.CreateUser(etudiant)
		if err != nil {
			if err.Error() == "user already exists" {
				fmt.Printf("Etudiant %s already exists\n", etudiant.Email)
			} else {
				log.Printf("Failed to create Etudiant %s: %v", etudiant.Email, err)
			}
		} else {
			fmt.Printf("Etudiant %s created successfully\n", etudiant.Email)
		}
	}

	fmt.Println("\nSeed completed successfully!")
	fmt.Println("\nTest accounts:")
	fmt.Println("Super Admin: superadmin@eduqr.com / superadmin123")
	fmt.Println("Admin: admin@eduqr.com / admin123")
	fmt.Println("Professeur: prof1@eduqr.com / prof123")
	fmt.Println("Etudiant: etudiant1@eduqr.com / student123")
}
