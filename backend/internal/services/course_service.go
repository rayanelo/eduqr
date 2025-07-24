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
		return nil, err
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

	// Mettre à jour le cours
	if err := s.courseRepo.UpdateCourse(course); err != nil {
		return nil, err
	}

	// Récupérer le cours mis à jour avec ses relations
	updatedCourse, err := s.courseRepo.GetCourseByID(id)
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
