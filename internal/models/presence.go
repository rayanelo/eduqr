package models

import (
	"time"

	"gorm.io/gorm"
)

// Status constants for presence
const (
	StatusPresent = "present" // Présent
	StatusLate    = "late"    // En retard
	StatusAbsent  = "absent"  // Absent
)

// Presence represents a student's attendance record
type Presence struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	StudentID uint           `json:"student_id" gorm:"not null;index"`
	Student   User           `json:"student" gorm:"foreignKey:StudentID"`
	CourseID  uint           `json:"course_id" gorm:"not null;index"`
	Course    Course         `json:"course" gorm:"foreignKey:CourseID"`
	Status    string         `json:"status" gorm:"default:'absent';index"` // present, late, absent
	ScannedAt *time.Time     `json:"scanned_at"`                           // Heure du scan du QR code
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PresenceResponse pour l'API
type PresenceResponse struct {
	ID        uint           `json:"id"`
	Student   UserResponse   `json:"student"`
	Course    CourseResponse `json:"course"`
	Status    string         `json:"status"`
	ScannedAt *time.Time     `json:"scanned_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ScanQRRequest pour le scan d'un QR code
type ScanQRRequest struct {
	QRCodeData string `json:"qr_code_data" binding:"required"`
}

// QRCodeInfo représente les informations d'un QR code
type QRCodeInfo struct {
	CourseID    uint      `json:"course_id"`
	CourseName  string    `json:"course_name"`
	SubjectName string    `json:"subject_name"`
	TeacherName string    `json:"teacher_name"`
	RoomName    string    `json:"room_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	QRCodeData  string    `json:"qr_code_data"`
	IsValid     bool      `json:"is_valid"`
}

// PresenceStatsResponse pour les statistiques de présence
type PresenceStatsResponse struct {
	TotalStudents   int64   `json:"total_students"`
	PresentStudents int64   `json:"present_students"`
	LateStudents    int64   `json:"late_students"`
	AbsentStudents  int64   `json:"absent_students"`
	AttendanceRate  float64 `json:"attendance_rate"`
}

// ToPresenceResponse convertit un Presence en PresenceResponse
func (p *Presence) ToPresenceResponse() PresenceResponse {
	return PresenceResponse{
		ID:        p.ID,
		Student:   UserToUserResponse(p.Student),
		Course:    p.Course.ToCourseResponse(),
		Status:    p.Status,
		ScannedAt: p.ScannedAt,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
