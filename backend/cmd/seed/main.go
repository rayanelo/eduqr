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
	err = database.AutoMigrate(&models.User{}, &models.AuditLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	auditLogRepo := repositories.NewAuditLogRepository()

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)
	auditLogService := services.NewAuditLogService(auditLogRepo)

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

	// Create audit logs for testing
	log.Println("Creating audit logs...")

	// Get users for audit logs
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Println("Failed to get users for audit logs:", err)
	} else if len(users) > 0 {
		// Create some test audit logs
		emptyStr := ""
		oldValues1 := `{"title": "Mathématiques Avancées", "duration": 120}`
		newValues1 := `{"title": "Mathématiques Avancées", "duration": 90}`
		newValues2 := `{"name": "Salle A101", "capacity": 30, "type": "classroom"}`
		oldValues2 := `{"email": "prof2@eduqr.com", "role": "professeur"}`

		testLogs := []models.AuditLog{
			{
				UserID:       users[0].ID, // Super Admin
				UserEmail:    users[0].Email,
				UserRole:     users[0].Role,
				Action:       models.ActionLogin,
				ResourceType: "user",
				ResourceID:   &users[0].ID,
				Description:  "Connexion de l'utilisateur Super Admin",
				IPAddress:    "192.168.1.100",
				UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				OldValues:    &emptyStr,
				NewValues:    &emptyStr,
			},
			{
				UserID:       users[1].ID, // Admin
				UserEmail:    users[1].Email,
				UserRole:     users[1].Role,
				Action:       models.ActionCreate,
				ResourceType: "room",
				ResourceID:   nil,
				Description:  "Création d'une nouvelle salle de cours",
				IPAddress:    "192.168.1.101",
				UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
				OldValues:    &emptyStr,
				NewValues:    &newValues2,
			},
			{
				UserID:       users[2].ID, // Prof 1
				UserEmail:    users[2].Email,
				UserRole:     users[2].Role,
				Action:       models.ActionUpdate,
				ResourceType: "course",
				ResourceID:   nil,
				Description:  "Modification du cours de mathématiques",
				IPAddress:    "192.168.1.102",
				UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
				OldValues:    &oldValues1,
				NewValues:    &newValues1,
			},
			{
				UserID:       users[0].ID, // Super Admin
				UserEmail:    users[0].Email,
				UserRole:     users[0].Role,
				Action:       models.ActionDelete,
				ResourceType: "user",
				ResourceID:   &users[3].ID,
				Description:  "Suppression de l'utilisateur prof2@eduqr.com",
				IPAddress:    "192.168.1.100",
				UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				OldValues:    &oldValues2,
				NewValues:    &emptyStr,
			},
			{
				UserID:       users[1].ID, // Admin
				UserEmail:    users[1].Email,
				UserRole:     users[1].Role,
				Action:       models.ActionLogout,
				ResourceType: "user",
				ResourceID:   &users[1].ID,
				Description:  "Déconnexion de l'utilisateur Admin",
				IPAddress:    "192.168.1.101",
				UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
				OldValues:    &emptyStr,
				NewValues:    &emptyStr,
			},
		}

		// Add some time variation to the logs
		for i, auditLog := range testLogs {
			auditLog.CreatedAt = time.Now().Add(-time.Duration(i*2) * time.Hour)

			// Convert to AuditLogRequest
			req := &models.AuditLogRequest{
				UserID:       auditLog.UserID,
				UserEmail:    auditLog.UserEmail,
				UserRole:     auditLog.UserRole,
				Action:       auditLog.Action,
				ResourceType: auditLog.ResourceType,
				ResourceID:   auditLog.ResourceID,
				Description:  auditLog.Description,
				IPAddress:    auditLog.IPAddress,
				UserAgent:    auditLog.UserAgent,
				OldValues:    auditLog.OldValues,
				NewValues:    auditLog.NewValues,
			}

			if _, err := auditLogService.CreateAuditLog(req); err != nil {
				log.Printf("Failed to create audit log %d: %v", i+1, err)
			} else {
				log.Printf("Audit log %d created successfully", i+1)
			}
		}
	}

	log.Println("Seed completed successfully!")

	log.Println("\nTest accounts:")
	log.Println("Super Admin: superadmin@eduqr.com / superadmin123")
	log.Println("Admin: admin@eduqr.com / admin123")
	log.Println("Professeur: prof1@eduqr.com / prof123")
	log.Println("Etudiant: etudiant1@eduqr.com / student123")
}
