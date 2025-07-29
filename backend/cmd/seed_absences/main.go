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
	err = database.AutoMigrate(&models.User{}, &models.Subject{}, &models.Room{}, &models.Course{}, &models.Absence{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories and services
	db := database.GetDB()
	userRepo := repositories.NewUserRepository()
	courseRepo := repositories.NewCourseRepository(db)
	absenceRepo := repositories.NewAbsenceRepository(db)

	absenceService := services.NewAbsenceService(absenceRepo, courseRepo, userRepo)

	// Récupérer les utilisateurs existants
	users, err := userRepo.FindAll()
	if err != nil {
		log.Fatal("Failed to get users:", err)
	}

	// Récupérer les cours existants
	courses, err := courseRepo.GetAllCourses()
	if err != nil {
		log.Fatal("Failed to get courses:", err)
	}

	if len(users) == 0 || len(courses) == 0 {
		log.Println("No users or courses found. Please run seed scripts first.")
		return
	}

	// Trouver un étudiant
	var student *models.User
	for _, user := range users {
		if user.Role == models.RoleEtudiant {
			student = &user
			break
		}
	}

	if student == nil {
		log.Println("No student found. Please create a student first.")
		return
	}

	// Trouver des cours passés
	var pastCourses []models.Course
	now := time.Now()
	for _, course := range courses {
		if course.EndTime.Before(now) {
			pastCourses = append(pastCourses, course)
		}
	}

	if len(pastCourses) == 0 {
		log.Println("No past courses found. Please create some past courses first.")
		return
	}

	log.Println("Creating test absences...")

	// Créer quelques absences de test
	testAbsences := []*models.CreateAbsenceRequest{
		{
			CourseID:      pastCourses[0].ID,
			Justification: "Maladie avec certificat médical",
			DocumentPath:  "/uploads/justificatifs/certificat_medical_1.pdf",
		},
		{
			CourseID:      pastCourses[0].ID,
			Justification: "Rendez-vous médical urgent",
			DocumentPath:  "/uploads/justificatifs/rdv_medical_1.pdf",
		},
	}

	for i, absenceReq := range testAbsences {
		// Vérifier si l'absence existe déjà
		exists, err := absenceRepo.CheckAbsenceExists(student.ID, absenceReq.CourseID)
		if err != nil {
			log.Printf("Error checking absence existence: %v", err)
			continue
		}

		if exists {
			log.Printf("Absence for course %d already exists, skipping...", absenceReq.CourseID)
			continue
		}

		// Créer l'absence
		absence, err := absenceService.CreateAbsence(absenceReq, student.ID)
		if err != nil {
			log.Printf("Failed to create absence %d: %v", i+1, err)
		} else {
			fmt.Printf("Absence %d created successfully: ID=%d, Course=%s, Status=%s\n",
				i+1, absence.ID, absence.Course.Name, absence.Status)
		}
	}

	// Créer quelques absences avec différents statuts pour tester
	if len(pastCourses) > 1 {
		// Créer une absence approuvée
		approvedAbsence := &models.CreateAbsenceRequest{
			CourseID:      pastCourses[1].ID,
			Justification: "Absence justifiée pour événement familial",
			DocumentPath:  "/uploads/justificatifs/evenement_familial.pdf",
		}

		exists, err := absenceRepo.CheckAbsenceExists(student.ID, approvedAbsence.CourseID)
		if err == nil && !exists {
			absence, err := absenceService.CreateAbsence(approvedAbsence, student.ID)
			if err == nil {
				// Approuver l'absence
				reviewReq := &models.ReviewAbsenceRequest{
					Status:        models.StatusApproved,
					ReviewComment: "Justificatif accepté, absence justifiée",
				}

				// Trouver un professeur ou admin pour approuver
				var reviewer *models.User
				for _, user := range users {
					if user.Role == models.RoleProfesseur || user.Role == models.RoleAdmin {
						reviewer = &user
						break
					}
				}

				if reviewer != nil {
					_, err = absenceService.ReviewAbsence(absence.ID, reviewReq, reviewer.ID, reviewer.Role)
					if err != nil {
						log.Printf("Failed to approve absence: %v", err)
					} else {
						fmt.Printf("Absence approved successfully: ID=%d\n", absence.ID)
					}
				}
			}
		}
	}

	log.Println("Absence seeding completed!")
}
