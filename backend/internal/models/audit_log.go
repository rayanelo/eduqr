package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ActionType constants for audit logging
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionLogin  = "login"
	ActionLogout = "logout"
)

// ResourceType constants for audit logging
const (
	ResourceUser    = "user"
	ResourceRoom    = "room"
	ResourceSubject = "subject"
	ResourceCourse  = "course"
	ResourceEvent   = "event"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	UserEmail    string         `json:"user_email" gorm:"not null"`
	UserRole     string         `json:"user_role" gorm:"not null"`
	Action       string         `json:"action" gorm:"not null;index"`        // create, update, delete, login, logout
	ResourceType string         `json:"resource_type" gorm:"not null;index"` // user, room, subject, course, event
	ResourceID   *uint          `json:"resource_id" gorm:"index"`            // ID of the affected resource (nullable for login/logout)
	Description  string         `json:"description" gorm:"not null"`         // Human readable description
	OldValues    *string        `json:"old_values" gorm:"type:json"`         // JSON string of old values (for updates)
	NewValues    *string        `json:"new_values" gorm:"type:json"`         // JSON string of new values (for creates/updates)
	IPAddress    string         `json:"ip_address" gorm:"not null"`
	UserAgent    string         `json:"user_agent"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// AuditLogRequest represents a request to create an audit log
type AuditLogRequest struct {
	UserID       uint        `json:"user_id" binding:"required"`
	UserEmail    string      `json:"user_email" binding:"required"`
	UserRole     string      `json:"user_role" binding:"required"`
	Action       string      `json:"action" binding:"required"`
	ResourceType string      `json:"resource_type" binding:"required"`
	ResourceID   *uint       `json:"resource_id"`
	Description  string      `json:"description" binding:"required"`
	OldValues    interface{} `json:"old_values"`
	NewValues    interface{} `json:"new_values"`
	IPAddress    string      `json:"ip_address" binding:"required"`
	UserAgent    string      `json:"user_agent"`
}

// AuditLogResponse represents an audit log response
type AuditLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	UserRole     string    `json:"user_role"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *uint     `json:"resource_id"`
	Description  string    `json:"description"`
	OldValues    *string   `json:"old_values"`
	NewValues    *string   `json:"new_values"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditLogFilter represents filters for querying audit logs
type AuditLogFilter struct {
	UserID       *uint      `json:"user_id"`
	Action       string     `json:"action"`
	ResourceType string     `json:"resource_type"`
	ResourceID   *uint      `json:"resource_id"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Search       string     `json:"search"` // Search in description, user_email, etc.
	Page         int        `json:"page" binding:"min=1"`
	Limit        int        `json:"limit" binding:"min=1,max=100"`
}

// AuditLogListResponse represents a paginated list of audit logs
type AuditLogListResponse struct {
	Logs       []AuditLogResponse `json:"logs"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// Helper function to convert interface{} to JSON string
func ToJSONString(data interface{}) *string {
	if data == nil {
		return nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	jsonStr := string(jsonBytes)
	return &jsonStr
}

// Helper function to create audit log description
func CreateAuditDescription(action, resourceType, resourceName string) string {
	switch action {
	case ActionCreate:
		return "Création de " + resourceType + " : " + resourceName
	case ActionUpdate:
		return "Modification de " + resourceType + " : " + resourceName
	case ActionDelete:
		return "Suppression de " + resourceType + " : " + resourceName
	case ActionLogin:
		return "Connexion utilisateur"
	case ActionLogout:
		return "Déconnexion utilisateur"
	default:
		return action + " sur " + resourceType + " : " + resourceName
	}
}
