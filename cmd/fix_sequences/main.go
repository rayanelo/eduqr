package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"log"
)

func main() {
	log.Println("🔧 Réinitialisation des séquences PostgreSQL...")

	// Charger la configuration
	cfg := config.LoadConfig()

	// Se connecter à la base de données
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données: %v", err)
	}
	defer database.CloseDB()

	db := database.GetDB()

	// Requêtes pour réinitialiser les séquences
	queries := []string{
		// Réinitialiser la séquence des utilisateurs
		"SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM users), false);",

		// Réinitialiser la séquence des salles
		"SELECT setval('rooms_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM rooms), false);",

		// Réinitialiser la séquence des matières
		"SELECT setval('subjects_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM subjects), false);",

		// Réinitialiser la séquence des cours
		"SELECT setval('courses_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM courses), false);",

		// Réinitialiser la séquence des absences
		"SELECT setval('absences_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM absences), false);",

		// Réinitialiser la séquence des présences
		"SELECT setval('presences_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM presences), false);",

		// Réinitialiser la séquence des logs d'audit
		"SELECT setval('audit_logs_id_seq', (SELECT COALESCE(MAX(id), 0) + 1 FROM audit_logs), false);",
	}

	// Exécuter chaque requête
	for i, query := range queries {
		log.Printf("📝 Exécution de la requête %d/%d...", i+1, len(queries))

		var result int64
		if err := db.Raw(query).Scan(&result).Error; err != nil {
			log.Printf("⚠️  Erreur lors de l'exécution de la requête %d: %v", i+1, err)
			log.Printf("   Requête: %s", query)
		} else {
			log.Printf("✅ Séquence réinitialisée, prochain ID: %d", result)
		}
	}

	log.Println("🎉 Réinitialisation des séquences terminée !")

	// Afficher les informations sur les tables
	log.Println("\n📊 Informations sur les tables:")

	tables := []string{"users", "rooms", "subjects", "courses", "absences", "presences", "audit_logs"}

	for _, table := range tables {
		var count int64
		var maxID int64

		// Compter les enregistrements
		if err := db.Table(table).Count(&count).Error; err != nil {
			log.Printf("❌ Erreur lors du comptage de %s: %v", table, err)
			continue
		}

		// Obtenir l'ID maximum
		if err := db.Table(table).Select("COALESCE(MAX(id), 0)").Scan(&maxID).Error; err != nil {
			log.Printf("❌ Erreur lors de la récupération de l'ID max de %s: %v", table, err)
			continue
		}

		log.Printf("   %s: %d enregistrements, ID max: %d", table, count, maxID)
	}
}
