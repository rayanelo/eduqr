package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Role constants
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleProfesseur = "professeur"
	RoleEtudiant   = "etudiant"
)

// Role hierarchy - higher index means higher privileges
var RoleHierarchy = map[string]int{
	RoleSuperAdmin: 4,
	RoleAdmin:      3,
	RoleProfesseur: 2,
	RoleEtudiant:   1,
}

// ValidRoles contains all valid role values
var ValidRoles = []string{RoleSuperAdmin, RoleAdmin, RoleProfesseur, RoleEtudiant}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// CanManageRole checks if a user with role1 can manage a user with role2
func CanManageRole(managerRole, targetRole string) bool {
	// Super Admin can manage everyone
	if managerRole == RoleSuperAdmin {
		return true
	}

	// Admin can manage Professeur and Etudiant, but not other Admins or Super Admin
	if managerRole == RoleAdmin {
		return targetRole == RoleProfesseur || targetRole == RoleEtudiant
	}

	// Professeur and Etudiant cannot manage anyone
	return false
}

// CanViewRole checks if a user with role1 can view users with role2
func CanViewRole(viewerRole, targetRole string) bool {
	// Super Admin can view everyone
	if viewerRole == RoleSuperAdmin {
		return true
	}

	// Admin can view Professeur and Etudiant, but not other Admins (except self)
	if viewerRole == RoleAdmin {
		return targetRole == RoleProfesseur || targetRole == RoleEtudiant
	}

	// Professeur can view other Professeurs and Etudiants
	if viewerRole == RoleProfesseur {
		return targetRole == RoleProfesseur || targetRole == RoleEtudiant
	}

	// Etudiant can only view other Etudiants
	if viewerRole == RoleEtudiant {
		return targetRole == RoleEtudiant
	}

	return false
}

// GetViewableFields returns the fields that a user with the given role can view
func GetViewableFields(viewerRole, targetRole string) []string {
	// Super Admin can see all fields
	if viewerRole == RoleSuperAdmin {
		return []string{"id", "email", "first_name", "last_name", "phone", "address", "avatar", "role", "created_at", "updated_at"}
	}

	// Admin can see all fields for Professeur and Etudiant
	if viewerRole == RoleAdmin && (targetRole == RoleProfesseur || targetRole == RoleEtudiant) {
		return []string{"id", "email", "first_name", "last_name", "phone", "address", "avatar", "role", "created_at", "updated_at"}
	}

	// Professeur can see limited fields for other Professeurs and Etudiants
	if viewerRole == RoleProfesseur && (targetRole == RoleProfesseur || targetRole == RoleEtudiant) {
		return []string{"id", "first_name", "last_name", "role", "created_at"}
	}

	// Etudiant can only see name fields for other Etudiants
	if viewerRole == RoleEtudiant && targetRole == RoleEtudiant {
		return []string{"id", "first_name", "last_name"}
	}

	return []string{}
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	Avatar    string         `json:"avatar" gorm:"default:'/assets/images/avatars/default-avatar.png'"`
	Role      string         `json:"role" gorm:"default:'etudiant'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Address         string `json:"address" binding:"required"`
}

type CreateUserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Address         string `json:"address" binding:"required"`
	Role            string `json:"role" binding:"required"`
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	if !IsValidRole(r.Role) {
		return errors.New("invalid role")
	}
	return nil
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}
