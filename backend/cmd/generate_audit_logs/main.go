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

	log.Println("📝 Génération de logs d'audit de test...")

	// Récupérer quelques utilisateurs pour créer des logs
	var users []models.User
	err = db.Limit(5).Find(&users).Error
	if err != nil {
		log.Fatal("Erreur lors de la récupération des utilisateurs:", err)
	}

	if len(users) == 0 {
		log.Fatal("Aucun utilisateur trouvé pour créer des logs d'audit")
	}

	// auditRepo := repositories.NewAuditLogRepository()

	// Actions de test
	actions := []string{"login", "logout", "create", "update", "delete"}
	resourceTypes := []string{"user", "course", "room", "subject", "absence"}
	descriptions := []string{
		"Connexion utilisateur",
		"Déconnexion utilisateur",
		"Création d'un nouvel élément",
		"Modification d'un élément",
		"Suppression d'un élément",
	}

	logsCreated := 0

	// Créer des logs d'audit pour les 7 derniers jours
	for i := 0; i < 20; i++ {
		user := users[i%len(users)]
		action := actions[i%len(actions)]
		resourceType := resourceTypes[i%len(resourceTypes)]
		description := descriptions[i%len(descriptions)]

		// Date aléatoire dans les 7 derniers jours
		daysAgo := i % 7
		createdAt := time.Now().AddDate(0, 0, -daysAgo)

		resourceID := uint(i + 1)
		oldValues := `"{}"`
		newValues := `"{\"test\": \"value\"}"`

		auditLog := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   &resourceID,
			Description:  description,
			OldValues:    &oldValues,
			NewValues:    &newValues,
			IPAddress:    "192.168.1.100",
			UserAgent:    "Mozilla/5.0 (Test Browser)",
			CreatedAt:    createdAt,
		}

		err := db.Create(auditLog).Error
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création du log %d: %v", i+1, err)
		} else {
			logsCreated++
			log.Printf("✅ Log créé: %s - %s - %s", user.Email, action, description)
		}
	}

	log.Printf("")
	log.Printf("🎉 %d logs d'audit créés avec succès!", logsCreated)
	log.Printf("")
	log.Printf("📊 Vous pouvez maintenant vérifier la page journal d'activité!")
}
