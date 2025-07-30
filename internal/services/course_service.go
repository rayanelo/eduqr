package services

import (
	"fmt"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
)

type CourseService struct {
	courseRepo  *repositories.CourseRepository
	subjectRepo *repositories.SubjectRepository
	userRepo    *repositories.UserRepository
	roomRepo    *repositories.RoomRepository
}

func NewCourseService(
	courseRepo *repositories.CourseRepository,
	subjectRepo *repositories.SubjectRepository,
	userRepo *repositories.UserRepository,
	roomRepo *repositories.RoomRepository,
) *CourseService {
	return &CourseService{
		courseRepo:  courseRepo,
		subjectRepo: subjectRepo,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
	}
}

// GetAllCourses récupère tous les cours
func (s *CourseService) GetAllCourses() ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetAllCourses()
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}

	return responses, nil
}

// GetCourseByID récupère un cours par son ID
func (s *CourseService) GetCourseByID(id uint) (*models.CourseResponse, error) {
	course, err := s.courseRepo.GetCourseByID(id)
	if err != nil {
		return nil, err
	}

	response := course.ToCourseResponse()
	return &response, nil
}

// CreateCourse crée un nouveau cours
func (s *CourseService) CreateCourse(req *models.CreateCourseRequest) (*models.CourseResponse, error) {
	// Vérifier que la matière existe
	_, err := s.subjectRepo.GetSubjectByID(req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("matière non trouvée")
	}

	// Vérifier que l'enseignant existe et est un professeur
	teacher, err := s.userRepo.FindByID(req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("enseignant non trouvé")
	}
	if teacher.Role != models.RoleProfesseur {
		return nil, fmt.Errorf("l'utilisateur sélectionné n'est pas un enseignant")
	}

	// Vérifier que la salle existe
	_, err = s.roomRepo.GetRoomByID(req.RoomID)
	if err != nil {
		return nil, fmt.Errorf("salle non trouvée")
	}

	// Vérifier que la date de fin de récurrence est après la date de début
	if req.IsRecurring && req.RecurrenceEndDate != nil {
		if req.RecurrenceEndDate.Before(req.StartTime) || req.RecurrenceEndDate.Equal(req.StartTime) {
			return nil, fmt.Errorf("la date de fin de récurrence doit être après la date de début")
		}
	}

	// Créer le cours
	course := &models.Course{
		Name:              req.Name,
		SubjectID:         req.SubjectID,
		TeacherID:         req.TeacherID,
		RoomID:            req.RoomID,
		StartTime:         req.StartTime,
		Duration:          req.Duration,
		Description:       req.Description,
		IsRecurring:       req.IsRecurring,
		RecurrencePattern: req.RecurrencePattern,
		RecurrenceEndDate: req.RecurrenceEndDate,
		ExcludeHolidays:   req.ExcludeHolidays,
	}

	// Si c'est un cours récurrent, générer les cours
	if req.IsRecurring {
		if err := s.courseRepo.CreateCourse(course); err != nil {
			return nil, err
		}

		// Générer les cours récurrents
		if err := s.courseRepo.GenerateRecurringCourses(course); err != nil {
			// Supprimer le cours parent si la génération échoue
			s.courseRepo.DeleteCourse(course.ID)
			return nil, fmt.Errorf("erreur lors de la génération des cours récurrents: %v", err)
		}
	} else {
		// Cours ponctuel
		if err := s.courseRepo.CreateCourse(course); err != nil {
			return nil, err
		}
	}

	// Récupérer le cours créé avec ses relations
	createdCourse, err := s.courseRepo.GetCourseByID(course.ID)
	if err != nil {
		return nil, err
	}

	response := createdCourse.ToCourseResponse()
	return &response, nil
}

// UpdateCourse met à jour un cours existant
func (s *CourseService) UpdateCourse(id uint, req *models.UpdateCourseRequest) (*models.CourseResponse, error) {
	// Récupérer le cours existant
	course, err := s.courseRepo.GetCourseByID(id)
	if err != nil {
		// Si le cours n'existe pas, retourner une erreur explicite
		return nil, fmt.Errorf("cours avec l'ID %d non trouvé", id)
	}

	// Vérifier que la matière existe si elle est modifiée
	if req.SubjectID != 0 {
		_, err := s.subjectRepo.GetSubjectByID(req.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("matière non trouvée")
		}
		course.SubjectID = req.SubjectID
	}

	// Vérifier que l'enseignant existe et est un professeur si il est modifié
	if req.TeacherID != 0 {
		teacher, err := s.userRepo.FindByID(req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("enseignant non trouvé")
		}
		if teacher.Role != models.RoleProfesseur {
			return nil, fmt.Errorf("l'utilisateur sélectionné n'est pas un enseignant")
		}
		course.TeacherID = req.TeacherID
	}

	// Vérifier que la salle existe si elle est modifiée
	if req.RoomID != 0 {
		_, err := s.roomRepo.GetRoomByID(req.RoomID)
		if err != nil {
			return nil, fmt.Errorf("salle non trouvée")
		}
		course.RoomID = req.RoomID
	}

	// Mettre à jour les autres champs
	if req.Name != "" {
		course.Name = req.Name
	}
	if !req.StartTime.IsZero() {
		course.StartTime = req.StartTime
	}
	if req.Duration != 0 {
		course.Duration = req.Duration
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.IsRecurring != course.IsRecurring {
		course.IsRecurring = req.IsRecurring
	}
	if req.RecurrencePattern != nil {
		course.RecurrencePattern = req.RecurrencePattern
	}
	if req.RecurrenceEndDate != nil {
		course.RecurrenceEndDate = req.RecurrenceEndDate
	}
	course.ExcludeHolidays = req.ExcludeHolidays

	// Si c'est un cours récurrent, supprimer tous les cours récurrents existants et les régénérer
	if course.IsRecurring {
		fmt.Printf("DEBUG: Cours récurrent détecté - ID: %d, RecurrenceID: %v, IsRecurring: %v\n", course.ID, course.RecurrenceID, course.IsRecurring)

		// Vérifier si c'est un cours parent récurrent (pas un cours enfant)
		if course.RecurrenceID == nil {
			fmt.Printf("DEBUG: Cours parent récurrent - Suppression et régénération\n")
			// C'est un cours parent récurrent, supprimer toute la série et régénérer
			if err := s.courseRepo.DeleteRecurringCourses(id); err != nil {
				return nil, fmt.Errorf("erreur lors de la suppression des cours récurrents: %v", err)
			}

			// Réinitialiser l'ID pour créer un nouveau cours
			course.ID = 0
			course.CreatedAt = time.Time{}
			course.UpdatedAt = time.Time{}

			// Créer le nouveau cours parent
			if err := s.courseRepo.CreateCourse(course); err != nil {
				return nil, err
			}

			// Régénérer les cours récurrents
			if err := s.courseRepo.GenerateRecurringCourses(course); err != nil {
				// Supprimer le cours parent si la génération échoue
				s.courseRepo.DeleteCourse(course.ID)
				return nil, fmt.Errorf("erreur lors de la régénération des cours récurrents: %v", err)
			}

			// Le cours a été recréé avec un nouvel ID, retourner directement la réponse
			response := course.ToCourseResponse()
			return &response, nil
		} else {
			fmt.Printf("DEBUG: Cours enfant récurrent - Modification interdite\n")
			// C'est un cours enfant d'une série récurrente, on ne peut pas le modifier directement
			return nil, fmt.Errorf("impossible de modifier un cours récurrent individuel. Modifiez le cours parent de la série.")
		}
	} else {
		fmt.Printf("DEBUG: Cours ponctuel - Modification normale\n")
		// Cours ponctuel
		if err := s.courseRepo.UpdateCourse(course); err != nil {
			return nil, err
		}
	}

	// Récupérer le cours mis à jour avec ses relations
	updatedCourse, err := s.courseRepo.GetCourseByID(course.ID)
	if err != nil {
		return nil, err
	}

	response := updatedCourse.ToCourseResponse()
	return &response, nil
}

// DeleteCourse supprime un cours
func (s *CourseService) DeleteCourse(id uint) error {
	// Récupérer le cours pour vérifier s'il est récurrent
	course, err := s.courseRepo.GetCourseByID(id)
	if err != nil {
		return err
	}

	// Si c'est un cours parent récurrent, supprimer toute la série
	if course.IsRecurring && course.RecurrenceID == nil {
		return s.courseRepo.DeleteRecurringCourses(id)
	}

	// Sinon, supprimer seulement ce cours
	return s.courseRepo.DeleteCourse(id)
}

// GetCoursesByDateRange récupère les cours dans une plage de dates
func (s *CourseService) GetCoursesByDateRange(startDate, endDate time.Time) ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetCoursesByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}

	return responses, nil
}

// GetCoursesByRoom récupère les cours d'une salle
func (s *CourseService) GetCoursesByRoom(roomID uint) ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetCoursesByRoom(roomID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}

	return responses, nil
}

// GetCoursesByRoomAndDate récupère les cours d'une salle pour une date spécifique
func (s *CourseService) GetCoursesByRoomAndDate(roomID uint, targetDate time.Time) ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetCoursesByRoomAndDate(roomID, targetDate)
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}

	return responses, nil
}

// GetCoursesByTeacher récupère les cours d'un enseignant
func (s *CourseService) GetCoursesByTeacher(teacherID uint) ([]models.CourseResponse, error) {
	courses, err := s.courseRepo.GetCoursesByTeacher(teacherID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}

	return responses, nil
}

// CheckConflicts vérifie les conflits pour un cours
func (s *CourseService) CheckConflicts(req *models.CreateCourseRequest) ([]models.ConflictInfo, error) {
	course := &models.Course{
		RoomID:    req.RoomID,
		StartTime: req.StartTime,
		Duration:  req.Duration,
	}

	return s.courseRepo.CheckConflicts(course)
}

// CheckConflictsForUpdate vérifie les conflits pour la modification d'un cours
func (s *CourseService) CheckConflictsForUpdate(courseID uint, req *models.UpdateCourseRequest) ([]models.ConflictInfo, error) {
	// Récupérer le cours existant
	existingCourse, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		// Si le cours n'existe pas, retourner une liste vide de conflits
		// car il ne peut pas y avoir de conflit avec un cours inexistant
		return []models.ConflictInfo{}, nil
	}

	// Créer un cours temporaire avec les nouvelles valeurs
	course := &models.Course{
		ID:        courseID,
		RoomID:    req.RoomID,
		StartTime: req.StartTime,
		Duration:  req.Duration,
	}

	// Utiliser les valeurs existantes si non modifiées
	if req.RoomID == 0 {
		course.RoomID = existingCourse.RoomID
	}
	if req.StartTime.IsZero() {
		course.StartTime = existingCourse.StartTime
	}
	if req.Duration == 0 {
		course.Duration = existingCourse.Duration
	}

	// Calculer l'heure de fin
	course.EndTime = course.StartTime.Add(time.Duration(course.Duration) * time.Minute)

	// Pour les cours récurrents, exclure le cours parent et tous ses enfants
	if existingCourse.IsRecurring && existingCourse.RecurrenceID == nil {
		// C'est un cours parent récurrent, exclure toute la série
		return s.courseRepo.CheckConflictsExcluding(courseID, course)
	} else if existingCourse.IsRecurring && existingCourse.RecurrenceID != nil {
		// C'est un cours enfant récurrent, exclure le parent et tous les enfants
		return s.courseRepo.CheckConflictsExcluding(*existingCourse.RecurrenceID, course)
	} else {
		// Cours ponctuel, exclure seulement ce cours
		return s.courseRepo.CheckConflictsExcluding(courseID, course)
	}
}
