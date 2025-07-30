package models

import (
	"time"

	"gorm.io/gorm"
)

// Subject représente une matière dans le système
type Subject struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Code        string         `json:"code" gorm:"size:20"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// SubjectResponse représente la réponse pour une matière
type SubjectResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSubjectRequest représente la requête de création d'une matière
type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// UpdateSubjectRequest représente la requête de modification d'une matière
type UpdateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ToSubjectResponse convertit un Subject en SubjectResponse
func (s *Subject) ToSubjectResponse() SubjectResponse {
	return SubjectResponse{
		ID:          s.ID,
		Name:        s.Name,
		Code:        s.Code,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
