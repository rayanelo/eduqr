package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"fmt"
	"log"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Connexion à la base de données
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}
	db := database.GetDB()

	// Auto-migration
	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Erreur de migration:", err)
	}

	// Récupérer tous les utilisateurs
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatal("Erreur lors de la récupération des utilisateurs:", err)
	}

	fmt.Printf("Nombre total d'utilisateurs: %d\n", len(users))
	fmt.Println("\nListe des utilisateurs:")
	fmt.Println("ID | Email | Prénom | Nom | Rôle")
	fmt.Println("---|-------|--------|-----|-----")

	for _, user := range users {
		fmt.Printf("%d | %s | %s | %s | %s\n",
			user.ID,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Role)
	}

	// Compter les professeurs
	var profCount int64
	if err := db.Model(&models.User{}).Where("role = ?", "professeur").Count(&profCount).Error; err != nil {
		log.Fatal("Erreur lors du comptage des professeurs:", err)
	}

	fmt.Printf("\nNombre de professeurs: %d\n", profCount)

	// Afficher les professeurs
	var professors []models.User
	if err := db.Where("role = ?", "professeur").Find(&professors).Error; err != nil {
		log.Fatal("Erreur lors de la récupération des professeurs:", err)
	}

	if len(professors) > 0 {
		fmt.Println("\nListe des professeurs:")
		for _, prof := range professors {
			fmt.Printf("- %s %s (%s)\n", prof.FirstName, prof.LastName, prof.Email)
		}
	} else {
		fmt.Println("\nAucun professeur trouvé dans la base de données.")
		fmt.Println("Vous devez créer des utilisateurs avec le rôle 'professeur'.")
	}
}
