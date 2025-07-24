package models

import (
	"time"

	"gorm.io/gorm"
)

// Course représente un cours ou événement pédagogique
type Course struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name" gorm:"not null"`
	SubjectID         uint           `json:"subject_id" gorm:"not null"`
	Subject           Subject        `json:"subject" gorm:"foreignKey:SubjectID"`
	TeacherID         uint           `json:"teacher_id" gorm:"not null"`
	Teacher           User           `json:"teacher" gorm:"foreignKey:TeacherID"`
	RoomID            uint           `json:"room_id" gorm:"not null"`
	Room              Room           `json:"room" gorm:"foreignKey:RoomID"`
	StartTime         time.Time      `json:"start_time" gorm:"not null"`
	EndTime           time.Time      `json:"end_time" gorm:"not null"`
	Duration          int            `json:"duration" gorm:"not null"` // en minutes
	Description       string         `json:"description"`
	IsRecurring       bool           `json:"is_recurring" gorm:"default:false"`
	RecurrenceID      *uint          `json:"recurrence_id"`      // ID du cours parent pour les récurrences
	RecurrencePattern *string        `json:"recurrence_pattern"` // JSON string pour les jours de répétition
	RecurrenceEndDate *time.Time     `json:"recurrence_end_date"`
	ExcludeHolidays   bool           `json:"exclude_holidays" gorm:"default:true"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// CourseResponse pour l'API
type CourseResponse struct {
	ID                uint            `json:"id"`
	Name              string          `json:"name"`
	Subject           SubjectResponse `json:"subject"`
	Teacher           UserResponse    `json:"teacher"`
	Room              RoomResponse    `json:"room"`
	StartTime         time.Time       `json:"start_time"`
	EndTime           time.Time       `json:"end_time"`
	Duration          int             `json:"duration"`
	Description       string          `json:"description"`
	IsRecurring       bool            `json:"is_recurring"`
	RecurrenceID      *uint           `json:"recurrence_id"`
	RecurrencePattern *string         `json:"recurrence_pattern"`
	RecurrenceEndDate *time.Time      `json:"recurrence_end_date"`
	ExcludeHolidays   bool            `json:"exclude_holidays"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// CreateCourseRequest pour la création d'un cours
type CreateCourseRequest struct {
	Name              string     `json:"name" binding:"required"`
	SubjectID         uint       `json:"subject_id" binding:"required"`
	TeacherID         uint       `json:"teacher_id" binding:"required"`
	RoomID            uint       `json:"room_id" binding:"required"`
	StartTime         time.Time  `json:"start_time" binding:"required"`
	Duration          int        `json:"duration" binding:"required,min=15,max=480"` // 15min à 8h
	Description       string     `json:"description"`
	IsRecurring       bool       `json:"is_recurring"`
	RecurrencePattern *string    `json:"recurrence_pattern"` // ["monday", "wednesday", "friday"]
	RecurrenceEndDate *time.Time `json:"recurrence_end_date"`
	ExcludeHolidays   bool       `json:"exclude_holidays"`
}

// UpdateCourseRequest pour la modification d'un cours
type UpdateCourseRequest struct {
	Name              string     `json:"name"`
	SubjectID         uint       `json:"subject_id"`
	TeacherID         uint       `json:"teacher_id"`
	RoomID            uint       `json:"room_id"`
	StartTime         time.Time  `json:"start_time"`
	Duration          int        `json:"duration" binding:"min=15,max=480"`
	Description       string     `json:"description"`
	IsRecurring       bool       `json:"is_recurring"`
	RecurrencePattern *string    `json:"recurrence_pattern"`
	RecurrenceEndDate *time.Time `json:"recurrence_end_date"`
	ExcludeHolidays   bool       `json:"exclude_holidays"`
}

// RecurrencePattern représente les jours de répétition
type RecurrencePattern struct {
	Days []string `json:"days"` // ["monday", "tuesday", etc.]
}

// ConflictInfo pour les conflits de réservation
type ConflictInfo struct {
	Date       time.Time `json:"date"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	RoomName   string    `json:"room_name"`
	CourseName string    `json:"course_name"`
}

// ToCourseResponse convertit un Course en CourseResponse
func (c *Course) ToCourseResponse() CourseResponse {
	return CourseResponse{
		ID:                c.ID,
		Name:              c.Name,
		Subject:           c.Subject.ToSubjectResponse(),
		Teacher:           UserToUserResponse(c.Teacher),
		Room:              c.Room.ToRoomResponse(),
		StartTime:         c.StartTime,
		EndTime:           c.EndTime,
		Duration:          c.Duration,
		Description:       c.Description,
		IsRecurring:       c.IsRecurring,
		RecurrenceID:      c.RecurrenceID,
		RecurrencePattern: c.RecurrencePattern,
		RecurrenceEndDate: c.RecurrenceEndDate,
		ExcludeHolidays:   c.ExcludeHolidays,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

// UserToUserResponse convertit un User en UserResponse
func UserToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		ContactEmail: user.ContactEmail,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Address:      user.Address,
		Avatar:       user.Avatar,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
