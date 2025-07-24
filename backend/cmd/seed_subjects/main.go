package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"fmt"
	"log"
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
	err = database.AutoMigrate(&models.Subject{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories and services
	subjectRepo := repositories.NewSubjectRepository()
	subjectService := services.NewSubjectService(subjectRepo)

	// Create sample subjects
	subjects := []*models.CreateSubjectRequest{
		{
			Name:        "Mathématiques",
			Code:        "MATH101",
			Description: "Fondamentaux des mathématiques incluant l'algèbre, la géométrie et l'analyse",
		},
		{
			Name:        "Physique",
			Code:        "PHY101",
			Description: "Principes fondamentaux de la physique classique et moderne",
		},
		{
			Name:        "Chimie",
			Code:        "CHEM101",
			Description: "Étude de la composition, structure et propriétés de la matière",
		},
		{
			Name:        "Biologie",
			Code:        "BIO101",
			Description: "Science du vivant et des organismes",
		},
		{
			Name:        "Histoire",
			Code:        "HIST101",
			Description: "Étude des événements passés et de leur impact sur le présent",
		},
		{
			Name:        "Géographie",
			Code:        "GEO101",
			Description: "Étude de la Terre, de ses paysages et de ses populations",
		},
		{
			Name:        "Français",
			Code:        "FR101",
			Description: "Langue française, littérature et expression écrite et orale",
		},
		{
			Name:        "Anglais",
			Code:        "ENG101",
			Description: "Langue anglaise et communication internationale",
		},
		{
			Name:        "Informatique",
			Code:        "CS101",
			Description: "Programmation, algorithmes et technologies numériques",
		},
		{
			Name:        "Philosophie",
			Code:        "PHIL101",
			Description: "Réflexion sur les questions fondamentales de l'existence",
		},
	}

	for _, subject := range subjects {
		_, err = subjectService.CreateSubject(subject)
		if err != nil {
			if err.Error() == "une matière avec ce nom existe déjà" {
				fmt.Printf("Matière %s existe déjà\n", subject.Name)
			} else {
				log.Printf("Failed to create subject %s: %v", subject.Name, err)
			}
		} else {
			fmt.Printf("Matière %s créée avec succès\n", subject.Name)
		}
	}

	fmt.Println("\nSeed des matières terminé avec succès!")
	fmt.Println("\nMatières créées:")
	for _, subject := range subjects {
		fmt.Printf("- %s (%s): %s\n", subject.Name, subject.Code, subject.Description)
	}
}
