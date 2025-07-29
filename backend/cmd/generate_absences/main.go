package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
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

	db := database.GetDB()

	log.Println("📝 Génération d'absences de test...")

	// Récupérer quelques étudiants et cours
	var students []models.User
	err = db.Where("role = ?", "etudiant").Limit(5).Find(&students).Error
	if err != nil {
		log.Fatal("Erreur lors de la récupération des étudiants:", err)
	}

	var courses []models.Course
	err = db.Limit(3).Find(&courses).Error
	if err != nil {
		log.Fatal("Erreur lors de la récupération des cours:", err)
	}

	if len(students) == 0 {
		log.Fatal("Aucun étudiant trouvé pour créer des absences")
	}

	if len(courses) == 0 {
		log.Fatal("Aucun cours trouvé pour créer des absences")
	}

	// Statuts d'absence
	statuses := []string{"pending", "approved", "rejected"}
	reasons := []string{
		"Maladie",
		"Rendez-vous médical",
		"Problème de transport",
		"Événement familial",
		"Raison personnelle",
	}

	absencesCreated := 0

	// Créer des absences pour les 30 derniers jours
	for i := 0; i < 25; i++ {
		student := students[i%len(students)]
		course := courses[i%len(courses)]
		status := statuses[i%len(statuses)]
		reason := reasons[i%len(reasons)]

		// Date aléatoire dans les 30 derniers jours
		daysAgo := i % 30
		_ = time.Now().AddDate(0, 0, -daysAgo) // Pour référence future

		absence := &models.Absence{
			StudentID:     student.ID,
			CourseID:      course.ID,
			Justification: reason,
			Status:        status,
			ReviewComment: "Absence de test générée automatiquement",
			CreatedAt:     time.Now(),
		}

		err := db.Create(absence).Error
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de l'absence %d: %v", i+1, err)
		} else {
			absencesCreated++
			log.Printf("✅ Absence créée: %s - %s - %s - %s",
				student.Email, course.Name, status, reason)
		}
	}

	log.Printf("")
	log.Printf("🎉 %d absences créées avec succès!", absencesCreated)
	log.Printf("")
	log.Printf("📊 Répartition par statut:")

	// Afficher les statistiques
	var pendingCount, approvedCount, rejectedCount int64
	db.Model(&models.Absence{}).Where("status = ?", "pending").Count(&pendingCount)
	db.Model(&models.Absence{}).Where("status = ?", "approved").Count(&approvedCount)
	db.Model(&models.Absence{}).Where("status = ?", "rejected").Count(&rejectedCount)

	log.Printf("   - En attente: %d", pendingCount)
	log.Printf("   - Approuvées: %d", approvedCount)
	log.Printf("   - Rejetées: %d", rejectedCount)
	log.Printf("")
	log.Printf("📊 Vous pouvez maintenant vérifier la page gestion des absences!")
}
