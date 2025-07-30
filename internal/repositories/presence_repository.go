package repositories

import (
	"eduqr-backend/internal/models"

	"gorm.io/gorm"
)

type PresenceRepository struct {
	db *gorm.DB
}

func NewPresenceRepository(db *gorm.DB) *PresenceRepository {
	return &PresenceRepository{db: db}
}

// CreatePresence crée une nouvelle présence
func (r *PresenceRepository) CreatePresence(presence *models.Presence) error {
	return r.db.Create(presence).Error
}

// GetPresenceByID récupère une présence par ID
func (r *PresenceRepository) GetPresenceByID(id uint) (*models.Presence, error) {
	var presence models.Presence
	err := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").
		Where("id = ?", id).First(&presence).Error
	if err != nil {
		return nil, err
	}
	return &presence, nil
}

// GetPresenceByStudentAndCourse récupère une présence par étudiant et cours
func (r *PresenceRepository) GetPresenceByStudentAndCourse(studentID, courseID uint) (*models.Presence, error) {
	var presence models.Presence
	err := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").
		Where("student_id = ? AND course_id = ?", studentID, courseID).First(&presence).Error
	if err != nil {
		return nil, err
	}
	return &presence, nil
}

// GetPresencesByCourse récupère toutes les présences d'un cours
func (r *PresenceRepository) GetPresencesByCourse(courseID uint) ([]models.Presence, error) {
	var presences []models.Presence
	err := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").
		Where("course_id = ?", courseID).Find(&presences).Error
	return presences, err
}

// GetPresencesByStudent récupère toutes les présences d'un étudiant
func (r *PresenceRepository) GetPresencesByStudent(studentID uint) ([]models.Presence, error) {
	var presences []models.Presence
	err := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").
		Where("student_id = ?", studentID).Find(&presences).Error
	return presences, err
}

// UpdatePresence met à jour une présence
func (r *PresenceRepository) UpdatePresence(presence *models.Presence) error {
	return r.db.Save(presence).Error
}

// DeletePresence supprime une présence
func (r *PresenceRepository) DeletePresence(id uint) error {
	return r.db.Delete(&models.Presence{}, id).Error
}

// CheckPresenceExists vérifie si une présence existe déjà
func (r *PresenceRepository) CheckPresenceExists(studentID, courseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Presence{}).Where("student_id = ? AND course_id = ?", studentID, courseID).Count(&count).Error
	return count > 0, err
}

// GetPresenceStats récupère les statistiques de présence pour un cours
func (r *PresenceRepository) GetPresenceStats(courseID uint) (*models.PresenceStatsResponse, error) {
	var stats models.PresenceStatsResponse

	// Total des étudiants inscrits au cours (à implémenter selon votre logique)
	var totalStudents int64
	err := r.db.Model(&models.User{}).Where("role = ?", models.RoleEtudiant).Count(&totalStudents).Error
	if err != nil {
		return nil, err
	}
	stats.TotalStudents = totalStudents

	// Présents
	err = r.db.Model(&models.Presence{}).Where("course_id = ? AND status = ?", courseID, models.StatusPresent).Count(&stats.PresentStudents).Error
	if err != nil {
		return nil, err
	}

	// En retard
	err = r.db.Model(&models.Presence{}).Where("course_id = ? AND status = ?", courseID, models.StatusLate).Count(&stats.LateStudents).Error
	if err != nil {
		return nil, err
	}

	// Absents
	err = r.db.Model(&models.Presence{}).Where("course_id = ? AND status = ?", courseID, models.StatusAbsent).Count(&stats.AbsentStudents).Error
	if err != nil {
		return nil, err
	}

	// Calcul du taux de présence
	if totalStudents > 0 {
		stats.AttendanceRate = float64(stats.PresentStudents+stats.LateStudents) / float64(totalStudents) * 100
	}

	return &stats, nil
}

// GetPresencesWithFilters récupère les présences avec filtres
func (r *PresenceRepository) GetPresencesWithFilters(filters map[string]interface{}, page, limit int) ([]models.Presence, int64, error) {
	var presences []models.Presence
	var total int64

	query := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room")

	// Appliquer les filtres
	if courseID, ok := filters["course_id"].(uint); ok {
		query = query.Where("course_id = ?", courseID)
	}
	if studentID, ok := filters["student_id"].(uint); ok {
		query = query.Where("student_id = ?", studentID)
	}
	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}
	if startDate, ok := filters["start_date"].(string); ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(string); ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Compter le total
	err := query.Model(&models.Presence{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les données paginées
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&presences).Error

	return presences, total, err
}

// CreatePresenceForAllStudents crée des enregistrements de présence pour tous les étudiants d'un cours
func (r *PresenceRepository) CreatePresenceForAllStudents(courseID uint) error {
	// Récupérer tous les étudiants
	var students []models.User
	err := r.db.Where("role = ?", models.RoleEtudiant).Find(&students).Error
	if err != nil {
		return err
	}

	// Créer un enregistrement de présence pour chaque étudiant
	for _, student := range students {
		// Vérifier si la présence existe déjà
		exists, err := r.CheckPresenceExists(student.ID, courseID)
		if err != nil {
			continue
		}
		if exists {
			continue
		}

		presence := &models.Presence{
			StudentID: student.ID,
			CourseID:  courseID,
			Status:    models.StatusAbsent, // Par défaut absent
		}

		err = r.CreatePresence(presence)
		if err != nil {
			continue
		}
	}

	return nil
}
