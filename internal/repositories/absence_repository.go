package repositories

import (
	"time"

	"eduqr-backend/internal/models"

	"gorm.io/gorm"
)

type AbsenceRepository struct {
	db *gorm.DB
}

func NewAbsenceRepository(db *gorm.DB) *AbsenceRepository {
	return &AbsenceRepository{db: db}
}

// CreateAbsence crée une nouvelle absence
func (r *AbsenceRepository) CreateAbsence(absence *models.Absence) error {
	return r.db.Create(absence).Error
}

// GetAbsenceByID récupère une absence par son ID
func (r *AbsenceRepository) GetAbsenceByID(id uint) (*models.Absence, error) {
	var absence models.Absence
	err := r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").Preload("Reviewer").First(&absence, id).Error
	if err != nil {
		return nil, err
	}
	return &absence, nil
}

// GetAbsencesByStudent récupère toutes les absences d'un étudiant
func (r *AbsenceRepository) GetAbsencesByStudent(studentID uint, page, limit int) ([]models.Absence, int64, error) {
	var absences []models.Absence
	var total int64

	offset := (page - 1) * limit

	// Compter le total
	err := r.db.Model(&models.Absence{}).Where("student_id = ?", studentID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les absences avec pagination
	err = r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").Preload("Reviewer").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&absences).Error

	return absences, total, err
}

// GetAbsencesByTeacher récupère toutes les absences pour les cours d'un professeur
func (r *AbsenceRepository) GetAbsencesByTeacher(teacherID uint, page, limit int) ([]models.Absence, int64, error) {
	var absences []models.Absence
	var total int64

	offset := (page - 1) * limit

	// Compter le total
	err := r.db.Model(&models.Absence{}).
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les absences avec pagination
	err = r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").Preload("Reviewer").
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Order("absences.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&absences).Error

	return absences, total, err
}

// GetAllAbsences récupère toutes les absences (pour les admins)
func (r *AbsenceRepository) GetAllAbsences(page, limit int) ([]models.Absence, int64, error) {
	var absences []models.Absence
	var total int64

	offset := (page - 1) * limit

	// Compter le total
	err := r.db.Model(&models.Absence{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les absences avec pagination
	err = r.db.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&absences).Error

	return absences, total, err
}

// GetAbsencesWithFilters récupère les absences avec filtres
func (r *AbsenceRepository) GetAbsencesWithFilters(filters *models.AbsenceFilterRequest) ([]models.Absence, int64, error) {
	var absences []models.Absence
	var total int64

	query := r.db.Model(&models.Absence{})

	// Appliquer les filtres
	if filters.StudentID != nil {
		query = query.Where("student_id = ?", *filters.StudentID)
	}
	if filters.CourseID != nil {
		query = query.Where("course_id = ?", *filters.CourseID)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filters.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if filters.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filters.EndDate)
		if err == nil {
			query = query.Where("created_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	offset := (filters.Page - 1) * filters.Limit

	// Compter le total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les absences avec pagination
	err = query.Preload("Student").Preload("Course.Subject").Preload("Course.Teacher").Preload("Course.Room").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(filters.Limit).
		Find(&absences).Error

	return absences, total, err
}

// UpdateAbsence met à jour une absence
func (r *AbsenceRepository) UpdateAbsence(absence *models.Absence) error {
	return r.db.Save(absence).Error
}

// ReviewAbsence valide ou rejette une absence
func (r *AbsenceRepository) ReviewAbsence(id uint, status string, reviewerID uint, reviewComment string) error {
	now := time.Now()
	return r.db.Model(&models.Absence{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"reviewer_id":    reviewerID,
			"review_comment": reviewComment,
			"reviewed_at":    &now,
		}).Error
}

// DeleteAbsence supprime une absence (soft delete)
func (r *AbsenceRepository) DeleteAbsence(id uint) error {
	return r.db.Delete(&models.Absence{}, id).Error
}

// CheckAbsenceExists vérifie si une absence existe déjà pour un étudiant et un cours
func (r *AbsenceRepository) CheckAbsenceExists(studentID, courseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Absence{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error
	return count > 0, err
}

// GetAbsenceStats récupère les statistiques des absences
func (r *AbsenceRepository) GetAbsenceStats() (*models.AbsenceStatsResponse, error) {
	var stats models.AbsenceStatsResponse

	// Total des absences
	err := r.db.Model(&models.Absence{}).Count(&stats.TotalAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences en attente
	err = r.db.Model(&models.Absence{}).Where("status = ?", models.StatusPending).Count(&stats.PendingAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences approuvées
	err = r.db.Model(&models.Absence{}).Where("status = ?", models.StatusApproved).Count(&stats.ApprovedAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences rejetées
	err = r.db.Model(&models.Absence{}).Where("status = ?", models.StatusRejected).Count(&stats.RejectedAbsences).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAbsenceStatsByTeacher récupère les statistiques des absences pour un professeur
func (r *AbsenceRepository) GetAbsenceStatsByTeacher(teacherID uint) (*models.AbsenceStatsResponse, error) {
	var stats models.AbsenceStatsResponse

	// Total des absences pour les cours du professeur
	err := r.db.Model(&models.Absence{}).
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Count(&stats.TotalAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences en attente
	err = r.db.Model(&models.Absence{}).
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ? AND absences.status = ?", teacherID, models.StatusPending).
		Count(&stats.PendingAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences approuvées
	err = r.db.Model(&models.Absence{}).
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ? AND absences.status = ?", teacherID, models.StatusApproved).
		Count(&stats.ApprovedAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences rejetées
	err = r.db.Model(&models.Absence{}).
		Joins("JOIN courses ON absences.course_id = courses.id").
		Where("courses.teacher_id = ? AND absences.status = ?", teacherID, models.StatusRejected).
		Count(&stats.RejectedAbsences).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAbsenceStatsByStudent récupère les statistiques des absences pour un étudiant
func (r *AbsenceRepository) GetAbsenceStatsByStudent(studentID uint) (*models.AbsenceStatsResponse, error) {
	var stats models.AbsenceStatsResponse

	// Total des absences de l'étudiant
	err := r.db.Model(&models.Absence{}).Where("student_id = ?", studentID).Count(&stats.TotalAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences en attente
	err = r.db.Model(&models.Absence{}).Where("student_id = ? AND status = ?", studentID, models.StatusPending).Count(&stats.PendingAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences approuvées
	err = r.db.Model(&models.Absence{}).Where("student_id = ? AND status = ?", studentID, models.StatusApproved).Count(&stats.ApprovedAbsences).Error
	if err != nil {
		return nil, err
	}

	// Absences rejetées
	err = r.db.Model(&models.Absence{}).Where("student_id = ? AND status = ?", studentID, models.StatusRejected).Count(&stats.RejectedAbsences).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
