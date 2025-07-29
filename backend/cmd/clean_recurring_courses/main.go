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

	log.Println("🧹 Nettoyage des cours récurrents problématiques...")

	// Récupérer les cours parents récurrents avec des dates de fin égales aux dates de début
	var problematicCourses []models.Course
	err = db.Where("is_recurring = ? AND recurrence_end_date = start_time", true).Find(&problematicCourses).Error
	if err != nil {
		log.Fatal("Erreur lors de la récupération des cours problématiques:", err)
	}

	log.Printf("📊 Nombre de cours récurrents problématiques trouvés: %d", len(problematicCourses))

	for _, course := range problematicCourses {
		log.Printf("🗑️  Suppression du cours récurrent ID %d: %s", course.ID, course.Name)
		log.Printf("   - Date début: %s", course.StartTime)
		log.Printf("   - Date fin: %s", *course.RecurrenceEndDate)

		// Supprimer le cours parent et tous ses enfants
		if err := db.Delete(&models.Course{}, course.ID).Error; err != nil {
			log.Printf("⚠️  Erreur lors de la suppression du cours parent %d: %v", course.ID, err)
			continue
		}

		// Supprimer tous les cours enfants
		if err := db.Where("recurrence_id = ?", course.ID).Delete(&models.Course{}).Error; err != nil {
			log.Printf("⚠️  Erreur lors de la suppression des cours enfants du parent %d: %v", course.ID, err)
		}

		log.Printf("✅ Cours récurrent ID %d supprimé avec succès", course.ID)
	}

	log.Printf("")
	log.Printf("🎉 Nettoyage terminé!")
}
