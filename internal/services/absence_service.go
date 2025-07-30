package services

import (
	"fmt"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
)

type AbsenceService struct {
	absenceRepo *repositories.AbsenceRepository
	courseRepo  *repositories.CourseRepository
	userRepo    *repositories.UserRepository
}

func NewAbsenceService(
	absenceRepo *repositories.AbsenceRepository,
	courseRepo *repositories.CourseRepository,
	userRepo *repositories.UserRepository,
) *AbsenceService {
	return &AbsenceService{
		absenceRepo: absenceRepo,
		courseRepo:  courseRepo,
		userRepo:    userRepo,
	}
}

// CreateAbsence crée une nouvelle absence
func (s *AbsenceService) CreateAbsence(req *models.CreateAbsenceRequest, studentID uint) (*models.AbsenceResponse, error) {
	// Vérifier que l'étudiant existe et est bien un étudiant
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return nil, fmt.Errorf("étudiant non trouvé")
	}
	if student.Role != models.RoleEtudiant {
		return nil, fmt.Errorf("seuls les étudiants peuvent créer des absences")
	}

	// Vérifier que le cours existe
	course, err := s.courseRepo.GetCourseByID(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("cours non trouvé")
	}

	// Vérifier que le cours est passé
	if course.StartTime.After(time.Now()) {
		return nil, fmt.Errorf("vous ne pouvez justifier qu'un cours déjà passé")
	}

	// Vérifier qu'il n'y a pas déjà une absence pour ce cours et cet étudiant
	exists, err := s.absenceRepo.CheckAbsenceExists(studentID, req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification de l'absence existante")
	}
	if exists {
		return nil, fmt.Errorf("une absence a déjà été déclarée pour ce cours")
	}

	// Créer l'absence
	absence := &models.Absence{
		StudentID:     studentID,
		CourseID:      req.CourseID,
		Justification: req.Justification,
		DocumentPath:  req.DocumentPath,
		Status:        models.StatusPending,
	}

	err = s.absenceRepo.CreateAbsence(absence)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'absence: %v", err)
	}

	// Récupérer l'absence créée avec ses relations
	createdAbsence, err := s.absenceRepo.GetAbsenceByID(absence.ID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'absence créée")
	}

	response := createdAbsence.ToAbsenceResponse()
	return &response, nil
}

// GetAbsenceByID récupère une absence par son ID
func (s *AbsenceService) GetAbsenceByID(id uint, userID uint, userRole string) (*models.AbsenceResponse, error) {
	absence, err := s.absenceRepo.GetAbsenceByID(id)
	if err != nil {
		return nil, fmt.Errorf("absence non trouvée")
	}

	// Vérifier les permissions selon le rôle
	if !s.canViewAbsence(userID, userRole, absence) {
		return nil, fmt.Errorf("permissions insuffisantes pour voir cette absence")
	}

	response := absence.ToAbsenceResponse()
	return &response, nil
}

// GetAbsencesByStudent récupère les absences d'un étudiant
func (s *AbsenceService) GetAbsencesByStudent(studentID uint, page, limit int, userID uint, userRole string) ([]models.AbsenceResponse, int64, error) {
	// Vérifier les permissions
	if !s.canViewStudentAbsences(userID, userRole, studentID) {
		return nil, 0, fmt.Errorf("permissions insuffisantes")
	}

	absences, total, err := s.absenceRepo.GetAbsencesByStudent(studentID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération des absences")
	}

	responses := make([]models.AbsenceResponse, len(absences))
	for i, absence := range absences {
		responses[i] = absence.ToAbsenceResponse()
	}

	return responses, total, nil
}

// GetAbsencesByTeacher récupère les absences pour les cours d'un professeur
func (s *AbsenceService) GetAbsencesByTeacher(teacherID uint, page, limit int, userID uint, userRole string) ([]models.AbsenceResponse, int64, error) {
	// Vérifier les permissions
	if !s.canViewTeacherAbsences(userID, userRole, teacherID) {
		return nil, 0, fmt.Errorf("permissions insuffisantes")
	}

	absences, total, err := s.absenceRepo.GetAbsencesByTeacher(teacherID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération des absences")
	}

	responses := make([]models.AbsenceResponse, len(absences))
	for i, absence := range absences {
		responses[i] = absence.ToAbsenceResponse()
	}

	return responses, total, nil
}

// GetAllAbsences récupère toutes les absences (pour les admins)
func (s *AbsenceService) GetAllAbsences(page, limit int, userRole string) ([]models.AbsenceResponse, int64, error) {
	// Vérifier les permissions
	if !s.isAdmin(userRole) {
		return nil, 0, fmt.Errorf("permissions insuffisantes")
	}

	absences, total, err := s.absenceRepo.GetAllAbsences(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération des absences")
	}

	responses := make([]models.AbsenceResponse, len(absences))
	for i, absence := range absences {
		responses[i] = absence.ToAbsenceResponse()
	}

	return responses, total, nil
}

// GetAbsencesWithFilters récupère les absences avec filtres
func (s *AbsenceService) GetAbsencesWithFilters(filters *models.AbsenceFilterRequest, userID uint, userRole string) ([]models.AbsenceResponse, int64, error) {
	// Vérifier les permissions selon le rôle
	if !s.canUseFilters(userRole) {
		return nil, 0, fmt.Errorf("permissions insuffisantes")
	}

	absences, total, err := s.absenceRepo.GetAbsencesWithFilters(filters)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération des absences")
	}

	responses := make([]models.AbsenceResponse, len(absences))
	for i, absence := range absences {
		responses[i] = absence.ToAbsenceResponse()
	}

	return responses, total, nil
}

// ReviewAbsence valide ou rejette une absence
func (s *AbsenceService) ReviewAbsence(id uint, req *models.ReviewAbsenceRequest, reviewerID uint, reviewerRole string) (*models.AbsenceResponse, error) {
	// Récupérer l'absence
	absence, err := s.absenceRepo.GetAbsenceByID(id)
	if err != nil {
		return nil, fmt.Errorf("absence non trouvée")
	}

	// Vérifier que l'absence n'a pas déjà été traitée
	if absence.Status != models.StatusPending {
		return nil, fmt.Errorf("cette absence a déjà été traitée")
	}

	// Vérifier les permissions pour traiter cette absence
	if !s.canReviewAbsence(reviewerID, reviewerRole, absence) {
		return nil, fmt.Errorf("permissions insuffisantes pour traiter cette absence")
	}

	// Traiter l'absence
	err = s.absenceRepo.ReviewAbsence(id, req.Status, reviewerID, req.ReviewComment)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du traitement de l'absence")
	}

	// Récupérer l'absence mise à jour
	updatedAbsence, err := s.absenceRepo.GetAbsenceByID(id)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'absence mise à jour")
	}

	response := updatedAbsence.ToAbsenceResponse()
	return &response, nil
}

// DeleteAbsence supprime une absence
func (s *AbsenceService) DeleteAbsence(id uint, userID uint, userRole string) error {
	// Récupérer l'absence
	absence, err := s.absenceRepo.GetAbsenceByID(id)
	if err != nil {
		return fmt.Errorf("absence non trouvée")
	}

	// Vérifier les permissions
	if !s.canDeleteAbsence(userID, userRole, absence) {
		return fmt.Errorf("permissions insuffisantes pour supprimer cette absence")
	}

	return s.absenceRepo.DeleteAbsence(id)
}

// GetAbsenceStats récupère les statistiques des absences
func (s *AbsenceService) GetAbsenceStats(userID uint, userRole string) (*models.AbsenceStatsResponse, error) {
	switch userRole {
	case models.RoleSuperAdmin, models.RoleAdmin:
		return s.absenceRepo.GetAbsenceStats()
	case models.RoleProfesseur:
		return s.absenceRepo.GetAbsenceStatsByTeacher(userID)
	case models.RoleEtudiant:
		return s.absenceRepo.GetAbsenceStatsByStudent(userID)
	default:
		return nil, fmt.Errorf("rôle non reconnu")
	}
}

// Méthodes de vérification des permissions

func (s *AbsenceService) canViewAbsence(userID uint, userRole string, absence *models.Absence) bool {
	// Super Admin et Admin peuvent voir toutes les absences
	if s.isAdmin(userRole) {
		return true
	}

	// Professeur peut voir les absences de ses cours
	if userRole == models.RoleProfesseur {
		return absence.Course.TeacherID == userID
	}

	// Étudiant peut voir ses propres absences
	if userRole == models.RoleEtudiant {
		return absence.StudentID == userID
	}

	return false
}

func (s *AbsenceService) canViewStudentAbsences(userID uint, userRole string, studentID uint) bool {
	// Super Admin et Admin peuvent voir toutes les absences
	if s.isAdmin(userRole) {
		return true
	}

	// Étudiant peut voir ses propres absences
	if userRole == models.RoleEtudiant {
		return userID == studentID
	}

	// Professeur peut voir les absences de ses étudiants (via les cours)
	// Cette logique est plus complexe et nécessiterait une vérification supplémentaire
	return userRole == models.RoleProfesseur
}

func (s *AbsenceService) canViewTeacherAbsences(userID uint, userRole string, teacherID uint) bool {
	// Super Admin et Admin peuvent voir toutes les absences
	if s.isAdmin(userRole) {
		return true
	}

	// Professeur peut voir les absences de ses propres cours
	if userRole == models.RoleProfesseur {
		return userID == teacherID
	}

	return false
}

func (s *AbsenceService) canReviewAbsence(reviewerID uint, reviewerRole string, absence *models.Absence) bool {
	// Super Admin et Admin peuvent traiter toutes les absences
	if s.isAdmin(reviewerRole) {
		return true
	}

	// Professeur peut traiter les absences de ses cours
	if reviewerRole == models.RoleProfesseur {
		return absence.Course.TeacherID == reviewerID
	}

	return false
}

func (s *AbsenceService) canDeleteAbsence(userID uint, userRole string, absence *models.Absence) bool {
	// Super Admin et Admin peuvent supprimer toutes les absences
	if s.isAdmin(userRole) {
		return true
	}

	// Étudiant peut supprimer ses propres absences si elles sont encore en attente
	if userRole == models.RoleEtudiant {
		return absence.StudentID == userID && absence.Status == models.StatusPending
	}

	return false
}

func (s *AbsenceService) canUseFilters(userRole string) bool {
	// Seuls les admins et professeurs peuvent utiliser les filtres avancés
	return s.isAdmin(userRole) || userRole == models.RoleProfesseur
}

func (s *AbsenceService) isAdmin(userRole string) bool {
	return userRole == models.RoleSuperAdmin || userRole == models.RoleAdmin
}
