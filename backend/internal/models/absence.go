package models

import (
	"time"

	"gorm.io/gorm"
)

// Status constants for absence justification
const (
	StatusPending  = "pending"  // En attente de validation
	StatusApproved = "approved" // Justificatif approuvé
	StatusRejected = "rejected" // Justificatif rejeté
)

// Absence represents a student absence with justification
type Absence struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	StudentID     uint           `json:"student_id" gorm:"not null;index"`
	Student       User           `json:"student" gorm:"foreignKey:StudentID"`
	CourseID      uint           `json:"course_id" gorm:"not null;index"`
	Course        Course         `json:"course" gorm:"foreignKey:CourseID"`
	Justification string         `json:"justification"`                         // Commentaire de l'étudiant
	DocumentPath  string         `json:"document_path"`                         // Chemin vers le fichier justificatif
	Status        string         `json:"status" gorm:"default:'pending';index"` // pending, approved, rejected
	ReviewerID    *uint          `json:"reviewer_id"`                           // ID de l'admin/professeur qui a validé/rejeté
	Reviewer      *User          `json:"reviewer" gorm:"foreignKey:ReviewerID"`
	ReviewComment string         `json:"review_comment"` // Commentaire du reviewer
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// AbsenceResponse pour l'API
type AbsenceResponse struct {
	ID            uint           `json:"id"`
	Student       UserResponse   `json:"student"`
	Course        CourseResponse `json:"course"`
	Justification string         `json:"justification"`
	DocumentPath  string         `json:"document_path"`
	Status        string         `json:"status"`
	Reviewer      *UserResponse  `json:"reviewer,omitempty"`
	ReviewComment string         `json:"review_comment"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// CreateAbsenceRequest pour la création d'une absence
type CreateAbsenceRequest struct {
	CourseID      uint   `json:"course_id" binding:"required"`
	Justification string `json:"justification"`
	DocumentPath  string `json:"document_path"`
}

// ReviewAbsenceRequest pour la validation/rejet d'une absence
type ReviewAbsenceRequest struct {
	Status        string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComment string `json:"review_comment" binding:"required"`
}

// AbsenceFilterRequest pour le filtrage des absences
type AbsenceFilterRequest struct {
	StudentID *uint   `json:"student_id"`
	CourseID  *uint   `json:"course_id"`
	Status    *string `json:"status"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Page      int     `json:"page" binding:"min=1"`
	Limit     int     `json:"limit" binding:"min=1,max=100"`
}

// AbsenceStatsResponse pour les statistiques des absences
type AbsenceStatsResponse struct {
	TotalAbsences    int64 `json:"total_absences"`
	PendingAbsences  int64 `json:"pending_absences"`
	ApprovedAbsences int64 `json:"approved_absences"`
	RejectedAbsences int64 `json:"rejected_absences"`
}

// ToAbsenceResponse convertit un Absence en AbsenceResponse
func (a *Absence) ToAbsenceResponse() AbsenceResponse {
	response := AbsenceResponse{
		ID:            a.ID,
		Student:       UserToUserResponse(a.Student),
		Course:        a.Course.ToCourseResponse(),
		Justification: a.Justification,
		DocumentPath:  a.DocumentPath,
		Status:        a.Status,
		ReviewComment: a.ReviewComment,
		ReviewedAt:    a.ReviewedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	if a.Reviewer != nil {
		reviewerResponse := UserToUserResponse(*a.Reviewer)
		response.Reviewer = &reviewerResponse
	}

	return response
}
