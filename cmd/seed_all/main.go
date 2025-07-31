package main

import (
	"fmt"
	"log"
	"time"

	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"eduqr-backend/pkg/utils"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Se connecter à la base de données
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// Auto-migrer les modèles
	if err := database.AutoMigrate(&models.User{}, &models.Room{}, &models.Subject{}, &models.Course{}, &models.Absence{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	fmt.Println("🚀 Démarrage du script de peuplement de la base de données...")

	// 1. Créer les utilisateurs
	fmt.Println("👥 Création des utilisateurs...")
	users := createUsers()
	fmt.Printf("✅ %d utilisateurs créés\n", len(users))

	// 2. Créer les salles
	fmt.Println("🏢 Création des salles...")
	rooms := createRooms()
	fmt.Printf("✅ %d salles créées\n", len(rooms))

	// 3. Créer les matières
	fmt.Println("📚 Création des matières...")
	subjects := createSubjects()
	fmt.Printf("✅ %d matières créées\n", len(subjects))

	// 4. Créer les cours
	fmt.Println("📅 Création des cours...")
	courses := createCourses(users, rooms, subjects)
	fmt.Printf("✅ %d cours créés\n", len(courses))

	// 5. Créer les absences
	fmt.Println("❌ Création des absences...")
	absences := createAbsences(users, courses)
	fmt.Printf("✅ %d absences créées\n", len(absences))

	fmt.Println("🎉 Peuplement de la base de données terminé avec succès!")
}

func createUsers() []models.User {
	users := []models.User{
		// Super Admin
		{
			Email:        "superadmin@eduqr.com",
			ContactEmail: "superadmin@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("SuperAdmin123!"); return p }(),
			FirstName:    "Super",
			LastName:     "Administrateur",
			Phone:        "+33123456789",
			Address:      "123 Rue de l'Administration, 75001 Paris",
			Role:         models.RoleSuperAdmin,
		},
		// Admin
		{
			Email:        "admin@eduqr.com",
			ContactEmail: "admin@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Admin123!"); return p }(),
			FirstName:    "Admin",
			LastName:     "Principal",
			Phone:        "+33123456790",
			Address:      "456 Rue de la Gestion, 75002 Paris",
			Role:         models.RoleAdmin,
		},
		// Professeurs
		{
			Email:        "prof1@eduqr.com",
			ContactEmail: "prof1@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Prof123!"); return p }(),
			FirstName:    "Marie",
			LastName:     "Dubois",
			Phone:        "+33123456791",
			Address:      "789 Rue des Professeurs, 75003 Paris",
			Role:         models.RoleProfesseur,
		},
		{
			Email:        "prof2@eduqr.com",
			ContactEmail: "prof2@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Prof123!"); return p }(),
			FirstName:    "Jean",
			LastName:     "Martin",
			Phone:        "+33123456792",
			Address:      "101 Rue de l'Enseignement, 75004 Paris",
			Role:         models.RoleProfesseur,
		},
		// Étudiants
		{
			Email:        "etudiant1@eduqr.com",
			ContactEmail: "etudiant1@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Etudiant123!"); return p }(),
			FirstName:    "Alice",
			LastName:     "Bernard",
			Phone:        "+33123456793",
			Address:      "202 Rue des Étudiants, 75005 Paris",
			Role:         models.RoleEtudiant,
		},
		{
			Email:        "etudiant2@eduqr.com",
			ContactEmail: "etudiant2@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Etudiant123!"); return p }(),
			FirstName:    "Thomas",
			LastName:     "Petit",
			Phone:        "+33123456794",
			Address:      "303 Avenue de l'Apprentissage, 75006 Paris",
			Role:         models.RoleEtudiant,
		},
		{
			Email:        "etudiant3@eduqr.com",
			ContactEmail: "etudiant3@eduqr.com",
			Password:     func() string { p, _ := utils.HashPassword("Etudiant123!"); return p }(),
			FirstName:    "Sophie",
			LastName:     "Roux",
			Phone:        "+33123456795",
			Address:      "404 Boulevard de l'Éducation, 75007 Paris",
			Role:         models.RoleEtudiant,
		},
	}

	for i := range users {
		if err := database.GetDB().Create(&users[i]).Error; err != nil {
			log.Printf("Erreur lors de la création de l'utilisateur %s: %v", users[i].Email, err)
		}
	}

	return users
}

func createRooms() []models.Room {
	// Créer d'abord la salle modulable parent
	modularRoom := models.Room{
		Name:      "Salle Modulable A",
		Building:  "Bâtiment Principal",
		Floor:     "1er étage",
		IsModular: true,
	}

	if err := database.GetDB().Create(&modularRoom).Error; err != nil {
		log.Printf("Erreur lors de la création de la salle modulable: %v", err)
	}

	// Créer les deux sous-salles
	subRooms := []models.Room{
		{
			Name:      "Salle Modulable A1",
			Building:  "Bâtiment Principal",
			Floor:     "1er étage",
			IsModular: false,
			ParentID:  &modularRoom.ID,
		},
		{
			Name:      "Salle Modulable A2",
			Building:  "Bâtiment Principal",
			Floor:     "1er étage",
			IsModular: false,
			ParentID:  &modularRoom.ID,
		},
	}

	for i := range subRooms {
		if err := database.GetDB().Create(&subRooms[i]).Error; err != nil {
			log.Printf("Erreur lors de la création de la sous-salle %s: %v", subRooms[i].Name, err)
		}
	}

	// Créer la salle normale
	normalRoom := models.Room{
		Name:      "Salle Normale B",
		Building:  "Bâtiment Principal",
		Floor:     "2ème étage",
		IsModular: false,
	}

	if err := database.GetDB().Create(&normalRoom).Error; err != nil {
		log.Printf("Erreur lors de la création de la salle normale: %v", err)
	}

	// Retourner toutes les salles
	rooms := []models.Room{modularRoom, subRooms[0], subRooms[1], normalRoom}
	return rooms
}

func createSubjects() []models.Subject {
	subjects := []models.Subject{
		{
			Name:        "Mathématiques",
			Code:        "MATH101",
			Description: "Cours de mathématiques fondamentales",
		},
		{
			Name:        "Informatique",
			Code:        "INFO101",
			Description: "Introduction à la programmation",
		},
		{
			Name:        "Physique",
			Code:        "PHYS101",
			Description: "Physique générale",
		},
		{
			Name:        "Anglais",
			Code:        "ANG101",
			Description: "Anglais technique",
		},
		{
			Name:        "Économie",
			Code:        "ECO101",
			Description: "Principes d'économie",
		},
	}

	for i := range subjects {
		if err := database.GetDB().Create(&subjects[i]).Error; err != nil {
			log.Printf("Erreur lors de la création de la matière %s: %v", subjects[i].Name, err)
		}
	}

	return subjects
}

func createCourses(users []models.User, rooms []models.Room, subjects []models.Subject) []models.Course {
	// Trouver les professeurs et les salles
	var profs []models.User
	var etudiants []models.User
	for _, user := range users {
		if user.Role == models.RoleProfesseur {
			profs = append(profs, user)
		} else if user.Role == models.RoleEtudiant {
			etudiants = append(etudiants, user)
		}
	}

	// Créer les cours récurrents (3 cours)
	recurringCourses := []models.Course{
		{
			Name:        "Mathématiques - Cours Récurrent",
			SubjectID:   subjects[0].ID, // Mathématiques
			TeacherID:   profs[0].ID,
			RoomID:      rooms[0].ID,                                  // Salle modulable
			StartTime:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), // Lundi 9h
			Duration:    120,                                          // 2h
			Description: "Cours de mathématiques récurrent chaque lundi",
			IsRecurring: true,
			RecurrencePattern: func() *string {
				pattern := `["monday"]`
				return &pattern
			}(),
			RecurrenceEndDate: func() *time.Time {
				end := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
				return &end
			}(),
		},
		{
			Name:        "Informatique - Cours Récurrent",
			SubjectID:   subjects[1].ID, // Informatique
			TeacherID:   profs[1].ID,
			RoomID:      rooms[1].ID,                                   // Sous-salle A1
			StartTime:   time.Date(2024, 1, 16, 14, 0, 0, 0, time.UTC), // Mardi 14h
			Duration:    180,                                           // 3h
			Description: "Cours d'informatique récurrent chaque mardi",
			IsRecurring: true,
			RecurrencePattern: func() *string {
				pattern := `["tuesday"]`
				return &pattern
			}(),
			RecurrenceEndDate: func() *time.Time {
				end := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
				return &end
			}(),
		},
		{
			Name:        "Physique - Cours Récurrent",
			SubjectID:   subjects[2].ID, // Physique
			TeacherID:   profs[0].ID,
			RoomID:      rooms[3].ID,                                   // Salle normale
			StartTime:   time.Date(2024, 1, 17, 16, 0, 0, 0, time.UTC), // Mercredi 16h
			Duration:    90,                                            // 1h30
			Description: "Cours de physique récurrent chaque mercredi",
			IsRecurring: true,
			RecurrencePattern: func() *string {
				pattern := `["wednesday"]`
				return &pattern
			}(),
			RecurrenceEndDate: func() *time.Time {
				end := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
				return &end
			}(),
		},
	}

	// Créer les cours normaux (2 cours)
	normalCourses := []models.Course{
		{
			Name:        "Anglais - Cours Unique",
			SubjectID:   subjects[3].ID, // Anglais
			TeacherID:   profs[1].ID,
			RoomID:      rooms[2].ID,                                   // Sous-salle A2
			StartTime:   time.Date(2024, 1, 18, 10, 0, 0, 0, time.UTC), // Jeudi 10h
			Duration:    60,                                            // 1h
			Description: "Cours d'anglais unique",
			IsRecurring: false,
		},
		{
			Name:        "Économie - Cours Unique",
			SubjectID:   subjects[4].ID, // Économie
			TeacherID:   profs[0].ID,
			RoomID:      rooms[0].ID,                                   // Salle modulable
			StartTime:   time.Date(2024, 1, 19, 15, 0, 0, 0, time.UTC), // Vendredi 15h
			Duration:    120,                                           // 2h
			Description: "Cours d'économie unique",
			IsRecurring: false,
		},
	}

	// Créer tous les cours
	allCourses := append(recurringCourses, normalCourses...)
	for i := range allCourses {
		// Calculer l'heure de fin
		allCourses[i].EndTime = allCourses[i].StartTime.Add(time.Duration(allCourses[i].Duration) * time.Minute)

		if err := database.GetDB().Create(&allCourses[i]).Error; err != nil {
			log.Printf("Erreur lors de la création du cours %s: %v", allCourses[i].Name, err)
		}
	}

	return allCourses
}

func createAbsences(users []models.User, courses []models.Course) []models.Absence {
	// Trouver les étudiants
	var etudiants []models.User
	for _, user := range users {
		if user.Role == models.RoleEtudiant {
			etudiants = append(etudiants, user)
		}
	}

	// Créer des absences avec différents statuts
	absences := []models.Absence{
		{
			StudentID:     etudiants[0].ID,
			CourseID:      courses[0].ID,
			Justification: "Maladie avec certificat médical",
			DocumentPath:  "/documents/certificat_medical_1.pdf",
			Status:        models.StatusApproved,
			ReviewerID:    &users[1].ID, // Admin
			ReviewComment: "Certificat médical valide",
			ReviewedAt:    func() *time.Time { t := time.Now(); return &t }(),
		},
		{
			StudentID:     etudiants[1].ID,
			CourseID:      courses[1].ID,
			Justification: "Problème de transport",
			DocumentPath:  "",
			Status:        models.StatusPending,
		},
		{
			StudentID:     etudiants[2].ID,
			CourseID:      courses[2].ID,
			Justification: "Rendez-vous médical",
			DocumentPath:  "/documents/rdv_medical_3.pdf",
			Status:        models.StatusRejected,
			ReviewerID:    &users[1].ID, // Admin
			ReviewComment: "Justificatif insuffisant",
			ReviewedAt:    func() *time.Time { t := time.Now(); return &t }(),
		},
		{
			StudentID:     etudiants[0].ID,
			CourseID:      courses[3].ID,
			Justification: "Problème familial urgent",
			DocumentPath:  "",
			Status:        models.StatusPending,
		},
		{
			StudentID:     etudiants[1].ID,
			CourseID:      courses[4].ID,
			Justification: "Absence justifiée par le médecin",
			DocumentPath:  "/documents/justificatif_medical_5.pdf",
			Status:        models.StatusApproved,
			ReviewerID:    &users[0].ID, // Super Admin
			ReviewComment: "Justificatif accepté",
			ReviewedAt:    func() *time.Time { t := time.Now(); return &t }(),
		},
	}

	for i := range absences {
		if err := database.GetDB().Create(&absences[i]).Error; err != nil {
			log.Printf("Erreur lors de la création de l'absence %d: %v", i+1, err)
		}
	}

	return absences
}
