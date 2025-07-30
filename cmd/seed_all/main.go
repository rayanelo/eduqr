package main

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("🚀 Démarrage du script de seed complet...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	log.Println("📊 Migration automatique des tables...")
	err = database.AutoMigrate(&models.User{}, &models.Subject{}, &models.Room{}, &models.Course{}, &models.Absence{}, &models.Presence{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories and services
	db := database.GetDB()
	userRepo := repositories.NewUserRepository()
	subjectRepo := repositories.NewSubjectRepository()
	roomRepo := repositories.NewRoomRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	absenceRepo := repositories.NewAbsenceRepository(db)
	presenceRepo := repositories.NewPresenceRepository(db)

	userService := services.NewUserService(userRepo, cfg.JWT.Secret, 24*time.Hour)
	subjectService := services.NewSubjectService(subjectRepo)
	roomService := services.NewRoomService(roomRepo)
	courseService := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)
	absenceService := services.NewAbsenceService(absenceRepo, courseRepo, userRepo)

	// 1. Création des utilisateurs
	log.Println("👥 Création des utilisateurs...")
	users := createUsers(userService)
	log.Printf("✅ %d utilisateurs créés", len(users))

	// 2. Création des matières
	log.Println("📚 Création des matières...")
	subjects := createSubjects(subjectService)
	log.Printf("✅ %d matières créées", len(subjects))

	// 3. Création des salles
	log.Println("🏫 Création des salles...")
	rooms := createRooms(roomService)
	log.Printf("✅ %d salles créées", len(rooms))

	// 4. Création des cours
	log.Println("📅 Création des cours...")
	courses := createCourses(courseService, users, subjects, rooms)
	log.Printf("✅ %d cours créés", len(courses))

	// 5. Création des présences
	log.Println("✅ Création des présences...")
	createPresences(presenceRepo, courses, users)
	log.Println("✅ Présences créées")

	// 6. Création des absences
	log.Println("❌ Création des absences...")
	createAbsences(absenceService, courses, users)
	log.Println("✅ Absences créées")

	log.Println("🎉 Script de seed terminé avec succès!")
}

func createUsers(userService *services.UserService) []models.User {
	userRequests := []*models.CreateUserRequest{
		{
			Email:           "superadmin@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Super",
			LastName:        "Admin",
			Phone:           "+1234567890",
			Address:         "123 Admin Street",
			Role:            "super_admin",
		},
		{
			Email:           "admin@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Admin",
			LastName:        "Principal",
			Phone:           "+1234567891",
			Address:         "456 Admin Avenue",
			Role:            "admin",
		},
		{
			Email:           "jean.dupont@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Jean",
			LastName:        "Dupont",
			Phone:           "+1234567892",
			Address:         "789 Teacher Street",
			Role:            "professeur",
		},
		{
			Email:           "marie.martin@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Marie",
			LastName:        "Martin",
			Phone:           "+1234567893",
			Address:         "321 Teacher Avenue",
			Role:            "professeur",
		},
		{
			Email:           "pierre.durand@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Pierre",
			LastName:        "Durand",
			Phone:           "+1234567894",
			Address:         "654 Teacher Road",
			Role:            "professeur",
		},
		{
			Email:           "alice.bernard@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Alice",
			LastName:        "Bernard",
			Phone:           "+1234567895",
			Address:         "987 Student Street",
			Role:            "etudiant",
		},
		{
			Email:           "bob.petit@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Bob",
			LastName:        "Petit",
			Phone:           "+1234567896",
			Address:         "147 Student Avenue",
			Role:            "etudiant",
		},
		{
			Email:           "claire.moreau@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Claire",
			LastName:        "Moreau",
			Phone:           "+1234567897",
			Address:         "258 Student Road",
			Role:            "etudiant",
		},
		{
			Email:           "david.leroy@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "David",
			LastName:        "Leroy",
			Phone:           "+1234567898",
			Address:         "369 Student Lane",
			Role:            "etudiant",
		},
		{
			Email:           "emma.roux@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Emma",
			LastName:        "Roux",
			Phone:           "+1234567899",
			Address:         "741 Student Drive",
			Role:            "etudiant",
		},
		{
			Email:           "francois.simon@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "François",
			LastName:        "Simon",
			Phone:           "+1234567800",
			Address:         "852 Student Court",
			Role:            "etudiant",
		},
		{
			Email:           "gabrielle.michel@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Gabrielle",
			LastName:        "Michel",
			Phone:           "+1234567801",
			Address:         "963 Student Place",
			Role:            "etudiant",
		},
		{
			Email:           "hugo.lefebvre@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Hugo",
			LastName:        "Lefebvre",
			Phone:           "+1234567802",
			Address:         "159 Student Way",
			Role:            "etudiant",
		},
		{
			Email:           "isabelle.leroy@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Isabelle",
			LastName:        "Leroy",
			Phone:           "+1234567803",
			Address:         "357 Student Circle",
			Role:            "etudiant",
		},
		{
			Email:           "julien.roux@eduqr.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			FirstName:       "Julien",
			LastName:        "Roux",
			Phone:           "+1234567804",
			Address:         "486 Student Square",
			Role:            "etudiant",
		},
	}

	var createdUsers []models.User
	for _, userReq := range userRequests {
		user, err := userService.CreateUser(userReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de l'utilisateur %s: %v", userReq.Email, err)
			continue
		}
		// Convertir UserResponse en User
		createdUser := models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
		}
		createdUsers = append(createdUsers, createdUser)
	}

	return createdUsers
}

func createSubjects(subjectService *services.SubjectService) []models.Subject {
	// D'abord, récupérer les matières existantes
	existingSubjects, err := subjectService.GetAllSubjects()
	if err != nil {
		log.Printf("⚠️ Erreur lors de la récupération des matières existantes: %v", err)
		return []models.Subject{}
	}

	// Si on a déjà des matières, les retourner
	if len(existingSubjects) > 0 {
		log.Println("📚 Utilisation des matières existantes...")
		var subjects []models.Subject
		for _, subject := range existingSubjects {
			subjects = append(subjects, models.Subject{
				ID:          subject.ID,
				Name:        subject.Name,
				Description: subject.Description,
				Code:        subject.Code,
			})
		}
		return subjects
	}

	// Sinon, créer de nouvelles matières
	subjectRequests := []*models.CreateSubjectRequest{
		{
			Name:        "Mathématiques",
			Description: "Cours de mathématiques avancées",
			Code:        "MATH101",
		},
		{
			Name:        "Physique",
			Description: "Cours de physique fondamentale",
			Code:        "PHYS101",
		},
		{
			Name:        "Informatique",
			Description: "Programmation et algorithmes",
			Code:        "INFO101",
		},
		{
			Name:        "Histoire",
			Description: "Histoire moderne et contemporaine",
			Code:        "HIST101",
		},
		{
			Name:        "Anglais",
			Description: "Cours d'anglais avancé",
			Code:        "ANG101",
		},
		{
			Name:        "Chimie",
			Description: "Chimie générale et organique",
			Code:        "CHIM101",
		},
	}

	var createdSubjects []models.Subject
	for _, subjectReq := range subjectRequests {
		subject, err := subjectService.CreateSubject(subjectReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la matière %s: %v", subjectReq.Name, err)
			continue
		}
		// Convertir SubjectResponse en Subject
		createdSubject := models.Subject{
			ID:          subject.ID,
			Name:        subject.Name,
			Description: subject.Description,
			Code:        subject.Code,
		}
		createdSubjects = append(createdSubjects, createdSubject)
	}

	return createdSubjects
}

func createRooms(roomService *services.RoomService) []models.Room {
	roomRequests := []*models.CreateRoomRequest{
		{
			Name:     "Salle A101",
			Building: "Bâtiment A",
			Floor:    "1er étage",
		},
		{
			Name:     "Salle A102",
			Building: "Bâtiment A",
			Floor:    "1er étage",
		},
		{
			Name:     "Salle B201",
			Building: "Bâtiment B",
			Floor:    "2ème étage",
		},
		{
			Name:     "Salle B202",
			Building: "Bâtiment B",
			Floor:    "2ème étage",
		},
		{
			Name:     "Labo Info C301",
			Building: "Bâtiment C",
			Floor:    "3ème étage",
		},
		{
			Name:     "Labo Physique C302",
			Building: "Bâtiment C",
			Floor:    "3ème étage",
		},
	}

	var createdRooms []models.Room
	for _, roomReq := range roomRequests {
		room, err := roomService.CreateRoom(roomReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création de la salle %s: %v", roomReq.Name, err)
			continue
		}
		// Convertir RoomResponse en Room
		createdRoom := models.Room{
			ID:       room.ID,
			Name:     room.Name,
			Building: room.Building,
			Floor:    room.Floor,
		}
		createdRooms = append(createdRooms, createdRoom)
	}

	return createdRooms
}

func createCourses(courseService *services.CourseService, users []models.User, subjects []models.Subject, rooms []models.Room) []models.Course {
	// Récupérer les professeurs et étudiants
	var teachers []models.User
	var students []models.User

	for _, user := range users {
		if user.Role == "professeur" {
			teachers = append(teachers, user)
		} else if user.Role == "etudiant" {
			students = append(students, user)
		}
	}

	// Créer des cours pour les 2 prochaines semaines
	now := time.Now()
	var courses []models.Course

	// Cours passés (pour les absences)
	for i := 1; i <= 5; i++ {
		startTime := now.AddDate(0, 0, -i*2) // Cours il y a 2, 4, 6, 8, 10 jours

		courseReq := &models.CreateCourseRequest{
			Name:        fmt.Sprintf("Cours %s - Session %d", subjects[i%len(subjects)].Name, i),
			Description: fmt.Sprintf("Description du cours %s", subjects[i%len(subjects)].Name),
			StartTime:   startTime,
			Duration:    120, // 2 heures en minutes
			TeacherID:   teachers[i%len(teachers)].ID,
			SubjectID:   subjects[i%len(subjects)].ID,
			RoomID:      rooms[i%len(rooms)].ID,
		}

		course, err := courseService.CreateCourse(courseReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création du cours %s: %v", courseReq.Name, err)
			continue
		}
		// Convertir CourseResponse en Course
		createdCourse := models.Course{
			ID:          course.ID,
			Name:        course.Name,
			Description: course.Description,
			StartTime:   course.StartTime,
			EndTime:     course.EndTime,
			TeacherID:   teachers[i%len(teachers)].ID,
			SubjectID:   subjects[i%len(subjects)].ID,
			RoomID:      rooms[i%len(rooms)].ID,
		}
		courses = append(courses, createdCourse)
	}

	// Cours futurs (pour les QR codes)
	for i := 1; i <= 10; i++ {
		startTime := now.AddDate(0, 0, i) // Cours dans 1, 2, 3... jours

		courseReq := &models.CreateCourseRequest{
			Name:        fmt.Sprintf("Cours %s - Session %d", subjects[i%len(subjects)].Name, i+5),
			Description: fmt.Sprintf("Description du cours %s", subjects[i%len(subjects)].Name),
			StartTime:   startTime,
			Duration:    120, // 2 heures en minutes
			TeacherID:   teachers[i%len(teachers)].ID,
			SubjectID:   subjects[i%len(subjects)].ID,
			RoomID:      rooms[i%len(rooms)].ID,
		}

		course, err := courseService.CreateCourse(courseReq)
		if err != nil {
			log.Printf("⚠️ Erreur lors de la création du cours %s: %v", courseReq.Name, err)
			continue
		}
		// Convertir CourseResponse en Course
		createdCourse := models.Course{
			ID:          course.ID,
			Name:        course.Name,
			Description: course.Description,
			StartTime:   course.StartTime,
			EndTime:     course.EndTime,
			TeacherID:   teachers[i%len(teachers)].ID,
			SubjectID:   subjects[i%len(subjects)].ID,
			RoomID:      rooms[i%len(rooms)].ID,
		}
		courses = append(courses, createdCourse)
	}

	return courses
}

func createPresences(presenceRepo *repositories.PresenceRepository, courses []models.Course, users []models.User) {
	// Récupérer les étudiants
	var students []models.User
	for _, user := range users {
		if user.Role == "etudiant" {
			students = append(students, user)
		}
	}

	// Créer des présences pour les cours passés
	for _, course := range courses {
		// Vérifier si le cours est passé
		if course.StartTime.After(time.Now()) {
			continue
		}

		for _, student := range students {
			// Simuler différents statuts de présence
			var status string
			var scannedAt time.Time

			// 70% de chance d'être présent, 20% en retard, 10% absent
			rand := time.Now().UnixNano() % 100
			if rand < 70 {
				status = "present"
				// Arrivé dans les 15 premières minutes
				scannedAt = course.StartTime.Add(time.Duration(rand%15) * time.Minute)
			} else if rand < 90 {
				status = "late"
				// Arrivé entre 15 et 30 minutes
				scannedAt = course.StartTime.Add(time.Duration(15+rand%15) * time.Minute)
			} else {
				status = "absent"
				scannedAt = time.Time{} // Pas de scan
			}

			// Créer la présence directement via le repository
			presence := &models.Presence{
				StudentID: student.ID,
				CourseID:  course.ID,
				Status:    status,
				ScannedAt: &scannedAt,
			}

			err := presenceRepo.CreatePresence(presence)
			if err != nil {
				log.Printf("⚠️ Erreur lors de la création de la présence pour l'étudiant %d dans le cours %d: %v", student.ID, course.ID, err)
			}
		}
	}
}

func createAbsences(absenceService *services.AbsenceService, courses []models.Course, users []models.User) {
	// Récupérer les étudiants
	var students []models.User
	for _, user := range users {
		if user.Role == "etudiant" {
			students = append(students, user)
		}
	}

	// Créer des absences pour les cours passés
	for _, course := range courses {
		// Vérifier si le cours est passé
		if course.StartTime.After(time.Now()) {
			continue
		}

		for _, student := range students {
			// 20% de chance d'avoir une absence justifiée
			if time.Now().UnixNano()%100 < 20 {
				absenceReq := &models.CreateAbsenceRequest{
					CourseID:      course.ID,
					Justification: "Absence justifiée pour raisons personnelles",
					DocumentPath:  "/uploads/justificatifs/absence_justifiee.pdf",
				}

				_, err := absenceService.CreateAbsence(absenceReq, student.ID)
				if err != nil {
					log.Printf("⚠️ Erreur lors de la création de l'absence pour l'étudiant %d dans le cours %d: %v", student.ID, course.ID, err)
				}
			}
		}
	}
}
