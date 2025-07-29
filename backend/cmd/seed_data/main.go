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
	err = database.AutoMigrate(&models.User{}, &models.Subject{}, &models.Room{}, &models.Course{}, &models.Presence{}, &models.Absence{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	subjectRepo := repositories.NewSubjectRepository()
	roomRepo := repositories.NewRoomRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	presenceRepo := repositories.NewPresenceRepository(db)
	absenceRepo := repositories.NewAbsenceRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)
	subjectService := services.NewSubjectService(subjectRepo)
	roomService := services.NewRoomService(roomRepo)
	courseService := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)
	presenceService := services.NewPresenceService(presenceRepo, courseRepo, userRepo)
	absenceService := services.NewAbsenceService(absenceRepo, courseRepo, userRepo)

	log.Println("🚀 Création des données de test...")

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

	var createdSubjects []models.Subject
	for _, subject := range subjects {
		createdSubject, err := subjectService.CreateSubject(&models.CreateSubjectRequest{
			Name:        subject.Name,
			Code:        subject.Code,
			Description: subject.Description,
		})
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la matière %s: %v", subject.Name, err)
		} else {
			// Convertir la réponse en modèle Subject
			subjectModel := models.Subject{
				ID:          createdSubject.ID,
				Name:        createdSubject.Name,
				Code:        createdSubject.Code,
				Description: createdSubject.Description,
			}
			createdSubjects = append(createdSubjects, subjectModel)
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

	var createdRooms []models.Room
	for _, room := range rooms {
		createdRoom, err := roomService.CreateRoom(&models.CreateRoomRequest{
			Name:      room.Name,
			Building:  room.Building,
			Floor:     room.Floor,
			IsModular: room.IsModular,
		})
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la salle %s: %v", room.Name, err)
		} else {
			// Convertir la réponse en modèle Room
			roomModel := models.Room{
				ID:        createdRoom.ID,
				Name:      createdRoom.Name,
				Building:  createdRoom.Building,
				Floor:     createdRoom.Floor,
				IsModular: createdRoom.IsModular,
			}
			createdRooms = append(createdRooms, roomModel)
			log.Printf("✅ Salle créée: %s", room.Name)
		}
	}

	// 3. Récupérer les utilisateurs existants
	log.Println("👥 Récupération des utilisateurs...")
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Fatal("Failed to get users:", err)
	}

	// Trouver les professeurs et étudiants
	var professors []models.User
	var students []models.User
	for _, user := range users {
		if user.Role == "professeur" {
			// Convertir UserResponse en User
			userModel := models.User{
				ID:    user.ID,
				Email: user.Email,
				Role:  user.Role,
			}
			professors = append(professors, userModel)
		} else if user.Role == "etudiant" {
			// Convertir UserResponse en User
			userModel := models.User{
				ID:    user.ID,
				Email: user.Email,
				Role:  user.Role,
			}
			students = append(students, userModel)
		}
	}

	log.Printf("📊 Utilisateurs trouvés: %d professeurs, %d étudiants", len(professors), len(students))

	// 4. Créer les cours
	log.Println("📅 Création des cours...")
	if len(professors) > 0 && len(createdSubjects) > 0 && len(createdRooms) > 0 {
		now := time.Now()
		weekly := "weekly"
		endDate := now.AddDate(0, 2, 0)

		// Cours pour cette semaine
		courses := []models.CreateCourseRequest{
			{
				Name:              "Mathématiques Avancées",
				SubjectID:         createdSubjects[0].ID, // Mathématiques
				TeacherID:         professors[0].ID,
				RoomID:            createdRooms[0].ID,
				StartTime:         now.AddDate(0, 0, 1).Add(9 * time.Hour), // Demain 9h
				Duration:          120,                                     // 2 heures
				Description:       "Cours de mathématiques avancées",
				IsRecurring:       true,
				RecurrencePattern: &weekly,
				RecurrenceEndDate: &endDate,
			},
			{
				Name:              "Physique Quantique",
				SubjectID:         createdSubjects[1].ID, // Physique
				TeacherID:         professors[0].ID,
				RoomID:            createdRooms[4].ID,                       // Labo
				StartTime:         now.AddDate(0, 0, 2).Add(14 * time.Hour), // Après-demain 14h
				Duration:          120,                                      // 2 heures
				Description:       "Introduction à la physique quantique",
				IsRecurring:       true,
				RecurrencePattern: &weekly,
				RecurrenceEndDate: &endDate,
			},
			{
				Name:              "Programmation Python",
				SubjectID:         createdSubjects[2].ID, // Informatique
				TeacherID:         professors[0].ID,
				RoomID:            createdRooms[4].ID,                       // Labo Info
				StartTime:         now.AddDate(0, 0, 3).Add(10 * time.Hour), // Dans 3 jours 10h
				Duration:          120,                                      // 2 heures
				Description:       "Apprentissage de Python",
				IsRecurring:       true,
				RecurrencePattern: &weekly,
				RecurrenceEndDate: &endDate,
			},
			{
				Name:              "Anglais Technique",
				SubjectID:         createdSubjects[3].ID, // Anglais
				TeacherID:         professors[0].ID,
				RoomID:            createdRooms[1].ID,
				StartTime:         now.AddDate(0, 0, 4).Add(13 * time.Hour), // Dans 4 jours 13h
				Duration:          90,                                       // 1h30
				Description:       "Anglais pour l'informatique",
				IsRecurring:       true,
				RecurrencePattern: &weekly,
				RecurrenceEndDate: &endDate,
			},
		}

		var createdCourses []models.Course
		for _, courseReq := range courses {
			createdCourse, err := courseService.CreateCourse(&courseReq)
			if err != nil {
				log.Printf("⚠️ Erreur lors de la création du cours %s: %v", courseReq.Name, err)
			} else {
				// Convertir CourseResponse en Course
				courseModel := models.Course{
					ID:   createdCourse.ID,
					Name: createdCourse.Name,
				}
				createdCourses = append(createdCourses, courseModel)
				log.Printf("✅ Cours créé: %s", courseReq.Name)
			}
		}

		// 5. Créer des présences pour les cours
		log.Println("✅ Création des présences...")
		for _, course := range createdCourses {
			// Créer des présences pour chaque étudiant
			for i, student := range students {
				if i >= 5 { // Limiter à 5 étudiants par cours pour les tests
					break
				}

				// Alterner entre présent et absent
				status := "present"
				if i%2 == 0 {
					status = "absent"
				}

				presenceReq := &models.CreatePresenceRequest{
					StudentID: student.ID,
					CourseID:  course.ID,
					Status:    status,
				}

				_, err := presenceService.CreatePresence(presenceReq)
				if err != nil {
					log.Printf("⚠️ Erreur lors de la création de la présence pour %s: %v", student.Email, err)
				} else {
					log.Printf("✅ Présence créée pour %s (cours: %s, statut: %s)", student.Email, course.Name, status)
				}
			}
		}

		// 6. Créer des absences justifiées
		log.Println("📝 Création des absences...")
		for i, student := range students {
			if i >= 3 { // Limiter à 3 absences pour les tests
				break
			}

			// Créer une absence pour le premier cours
			if len(createdCourses) > 0 {
				absenceReq := &models.CreateAbsenceRequest{
					StudentID:     student.ID,
					CourseID:      createdCourses[0].ID,
					Justification: "Maladie avec certificat médical",
					Status:        "pending", // En attente de validation
				}

				_, err := absenceService.CreateAbsence(absenceReq)
				if err != nil {
					log.Printf("⚠️ Erreur lors de la création de l'absence pour %s: %v", student.Email, err)
				} else {
					log.Printf("✅ Absence créée pour %s (cours: %s)", student.Email, createdCourses[0].Name)
				}
			}
		}

		log.Printf("🎉 Données de test créées avec succès!")
		log.Printf("📊 Résumé:")
		log.Printf("   - %d matières créées", len(createdSubjects))
		log.Printf("   - %d salles créées", len(createdRooms))
		log.Printf("   - %d cours créés", len(createdCourses))
		log.Printf("   - Présences et absences générées")
	} else {
		log.Println("❌ Impossible de créer les cours: manque de professeurs, matières ou salles")
	}
}
