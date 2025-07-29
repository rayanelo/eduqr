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
	err = database.AutoMigrate(&models.User{}, &models.Subject{}, &models.Room{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	subjectRepo := repositories.NewSubjectRepository()
	roomRepo := repositories.NewRoomRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)
	subjectService := services.NewSubjectService(subjectRepo)
	roomService := services.NewRoomService(roomRepo)

	log.Println("🚀 Création des données de test simples...")

	// 1. Créer les matières
	log.Println("📚 Création des matières...")
	subjects := []models.Subject{
		{Name: "Mathématiques", Code: "MATH101", Description: "Mathématiques fondamentales"},
		{Name: "Physique", Code: "PHYS101", Description: "Physique générale"},
		{Name: "Informatique", Code: "INFO101", Description: "Programmation de base"},
		{Name: "Anglais", Code: "ENG101", Description: "Anglais technique"},
		{Name: "Histoire", Code: "HIST101", Description: "Histoire contemporaine"},
		{Name: "Chimie", Code: "CHIM101", Description: "Chimie générale"},
		{Name: "Biologie", Code: "BIO101", Description: "Biologie cellulaire"},
		{Name: "Économie", Code: "ECO101", Description: "Principes d'économie"},
	}

	subjectCount := 0
	for _, subject := range subjects {
		_, err := subjectService.CreateSubject(&models.CreateSubjectRequest{
			Name:        subject.Name,
			Code:        subject.Code,
			Description: subject.Description,
		})
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la matière %s: %v", subject.Name, err)
		} else {
			subjectCount++
			log.Printf("✅ Matière créée: %s", subject.Name)
		}
	}

	// 2. Créer les salles
	log.Println("🏫 Création des salles...")
	rooms := []models.Room{
		{Name: "Salle A101", Building: "Bâtiment A", Floor: "1", IsModular: false},
		{Name: "Salle A102", Building: "Bâtiment A", Floor: "1", IsModular: false},
		{Name: "Salle B201", Building: "Bâtiment B", Floor: "2", IsModular: false},
		{Name: "Salle B202", Building: "Bâtiment B", Floor: "2", IsModular: false},
		{Name: "Labo Info C301", Building: "Bâtiment C", Floor: "3", IsModular: true},
		{Name: "Labo Physique C302", Building: "Bâtiment C", Floor: "3", IsModular: true},
		{Name: "Amphi Principal", Building: "Bâtiment A", Floor: "0", IsModular: false},
		{Name: "Salle de TD", Building: "Bâtiment B", Floor: "1", IsModular: false},
	}

	roomCount := 0
	for _, room := range rooms {
		_, err := roomService.CreateRoom(&models.CreateRoomRequest{
			Name:      room.Name,
			Building:  room.Building,
			Floor:     room.Floor,
			IsModular: room.IsModular,
		})
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la salle %s: %v", room.Name, err)
		} else {
			roomCount++
			log.Printf("✅ Salle créée: %s", room.Name)
		}
	}

	// 3. Afficher les utilisateurs existants
	log.Println("👥 Utilisateurs existants...")
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Printf("⚠️ Erreur lors de la récupération des utilisateurs: %v", err)
	} else {
		log.Printf("📊 %d utilisateurs trouvés", len(users))
		for _, user := range users {
			log.Printf("   - %s (%s) - %s", user.Email, user.FirstName+" "+user.LastName, user.Role)
		}
	}

	log.Printf("🎉 Données de test créées avec succès!")
	log.Printf("📊 Résumé:")
	log.Printf("   - %d matières créées", subjectCount)
	log.Printf("   - %d salles créées", roomCount)
	log.Printf("   - Utilisateurs existants listés")
	log.Printf("")
	log.Printf("✅ Vous pouvez maintenant lancer les scripts de test!")
}
