package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.GetDB()

	log.Println("🔍 Vérification des cours récurrents...")

	// Récupérer tous les cours
	var courses []models.Course
	err = db.Preload("Subject").Preload("Teacher").Preload("Room").Find(&courses).Error
	if err != nil {
		log.Fatal("Erreur lors de la récupération des cours:", err)
	}

	log.Printf("📊 Nombre total de cours: %d", len(courses))

	// Analyser les cours récurrents
	recurringParents := 0
	recurringChildren := 0

	for _, course := range courses {
		if course.IsRecurring {
			if course.RecurrenceID == nil {
				// C'est un cours parent récurrent
				recurringParents++
				log.Printf("🔄 Cours parent récurrent ID %d: %s", course.ID, course.Name)
				log.Printf("   - Pattern: %s", *course.RecurrencePattern)
				log.Printf("   - Date de fin: %s", *course.RecurrenceEndDate)
				log.Printf("   - Date de début: %s", course.StartTime)
				log.Printf("   - Durée: %d minutes", course.Duration)
			} else {
				// C'est un cours enfant récurrent
				recurringChildren++
				log.Printf("   📅 Cours enfant ID %d (parent: %d): %s", course.ID, *course.RecurrenceID, course.Name)
				log.Printf("      - Date: %s", course.StartTime)
			}
		}
	}

	log.Printf("")
	log.Printf("📈 Résumé des cours récurrents:")
	log.Printf("   - Cours parents récurrents: %d", recurringParents)
	log.Printf("   - Cours enfants récurrents: %d", recurringChildren)

	if recurringParents > 0 && recurringChildren == 0 {
		log.Printf("")
		log.Printf("⚠️  PROBLÈME DÉTECTÉ: Il y a %d cours parents récurrents mais aucun cours enfant!", recurringParents)
		log.Printf("   Cela indique que la génération des cours récurrents ne fonctionne pas.")
	} else if recurringParents > 0 && recurringChildren > 0 {
		log.Printf("")
		log.Printf("✅ Les cours récurrents semblent fonctionner correctement!")
	} else {
		log.Printf("")
		log.Printf("ℹ️  Aucun cours récurrent trouvé dans la base de données.")
	}
}
