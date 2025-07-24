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
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = database.AutoMigrate(&models.User{}, &models.Subject{}, &models.Room{}, &models.Course{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories and services
	db := database.GetDB()
	userRepo := repositories.NewUserRepository()
	subjectRepo := repositories.NewSubjectRepository()
	roomRepo := repositories.NewRoomRepository(db)
	courseRepo := repositories.NewCourseRepository(db)

	subjectService := services.NewSubjectService(subjectRepo)
	roomService := services.NewRoomService(roomRepo)
	courseService := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)

	// Créer des matières de test si elles n'existent pas
	subjects := []*models.CreateSubjectRequest{
		{Name: "Mathématiques", Code: "MATH", Description: "Cours de mathématiques"},
		{Name: "Physique", Code: "PHYS", Description: "Cours de physique"},
		{Name: "Informatique", Code: "INFO", Description: "Cours d'informatique"},
	}

	for _, subject := range subjects {
		_, err := subjectService.CreateSubject(subject)
		if err != nil && err.Error() != "subject already exists" {
			log.Printf("Failed to create subject %s: %v", subject.Name, err)
		} else {
			fmt.Printf("Subject %s created or already exists\n", subject.Name)
		}
	}

	// Créer des salles de test si elles n'existent pas
	rooms := []*models.CreateRoomRequest{
		{Name: "Salle A101", Building: "Bâtiment A", Floor: "1er étage", IsModular: false},
		{Name: "Salle B201", Building: "Bâtiment B", Floor: "2ème étage", IsModular: false},
		{Name: "Amphithéâtre 1", Building: "Bâtiment C", Floor: "Rez-de-chaussée", IsModular: false},
	}

	for _, room := range rooms {
		_, err := roomService.CreateRoom(room)
		if err != nil && err.Error() != "room already exists" {
			log.Printf("Failed to create room %s: %v", room.Name, err)
		} else {
			fmt.Printf("Room %s created or already exists\n", room.Name)
		}
	}

	// Récupérer les IDs des professeurs, matières et salles
	users, err := userRepo.FindAll()
	if err != nil {
		log.Fatal("Failed to get users:", err)
	}

	var prof1, prof2 *models.User
	for _, user := range users {
		if user.Role == "professeur" {
			if prof1 == nil {
				prof1 = &user
			} else if prof2 == nil {
				prof2 = &user
			}
		}
	}

	if prof1 == nil || prof2 == nil {
		log.Fatal("Need at least 2 professors")
	}

	subjectsList, err := subjectRepo.GetAllSubjects()
	if err != nil {
		log.Fatal("Failed to get subjects:", err)
	}

	roomsList, err := roomRepo.GetAllRooms(nil)
	if err != nil {
		log.Fatal("Failed to get rooms:", err)
	}

	if len(subjectsList) == 0 || len(roomsList) == 0 {
		log.Fatal("Need at least 1 subject and 1 room")
	}

	// Créer des cours récurrents de test
	recurringCourses := []*models.CreateCourseRequest{
		{
			Name:              "Mathématiques - Cours récurrent",
			SubjectID:         subjectsList[0].ID,
			TeacherID:         prof1.ID,
			RoomID:            roomsList[0].ID,
			StartTime:         time.Date(2024, 9, 2, 10, 0, 0, 0, time.UTC), // Lundi 2 septembre 2024 à 10h
			Duration:          120,                                          // 2 heures
			Description:       "Cours de mathématiques récurrent chaque lundi",
			IsRecurring:       true,
			RecurrencePattern: func() *string { s := `{"days": ["Monday"]}`; return &s }(),
			RecurrenceEndDate: func() *time.Time { t := time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC); return &t }(), // Jusqu'au 20 décembre
			ExcludeHolidays:   true,
		},
		{
			Name:              "Physique - TP récurrent",
			SubjectID:         subjectsList[1].ID,
			TeacherID:         prof2.ID,
			RoomID:            roomsList[1].ID,
			StartTime:         time.Date(2024, 9, 3, 14, 0, 0, 0, time.UTC), // Mardi 3 septembre 2024 à 14h
			Duration:          180,                                          // 3 heures
			Description:       "TP de physique récurrent chaque mardi et jeudi",
			IsRecurring:       true,
			RecurrencePattern: func() *string { s := `{"days": ["Tuesday", "Thursday"]}`; return &s }(),
			RecurrenceEndDate: func() *time.Time { t := time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC); return &t }(), // Jusqu'au 19 décembre
			ExcludeHolidays:   true,
		},
	}

	// Créer des cours ponctuels de test
	singleCourses := []*models.CreateCourseRequest{
		{
			Name:        "Informatique - Introduction",
			SubjectID:   subjectsList[2].ID,
			TeacherID:   prof1.ID,
			RoomID:      roomsList[2].ID,
			StartTime:   time.Now().Add(24 * time.Hour), // Demain
			Duration:    90,                             // 1h30
			Description: "Cours ponctuel d'introduction à l'informatique",
			IsRecurring: false,
		},
		{
			Name:        "Mathématiques - Révision",
			SubjectID:   subjectsList[0].ID,
			TeacherID:   prof2.ID,
			RoomID:      roomsList[0].ID,
			StartTime:   time.Now().Add(48 * time.Hour), // Après-demain
			Duration:    60,                             // 1 heure
			Description: "Cours ponctuel de révision",
			IsRecurring: false,
		},
	}

	// Créer les cours récurrents
	for _, course := range recurringCourses {
		_, err := courseService.CreateCourse(course)
		if err != nil {
			log.Printf("Failed to create recurring course %s: %v", course.Name, err)
		} else {
			fmt.Printf("Recurring course %s created successfully\n", course.Name)
		}
	}

	// Créer les cours ponctuels
	for _, course := range singleCourses {
		_, err := courseService.CreateCourse(course)
		if err != nil {
			log.Printf("Failed to create single course %s: %v", course.Name, err)
		} else {
			fmt.Printf("Single course %s created successfully\n", course.Name)
		}
	}

	fmt.Println("Course seeding completed!")
}
