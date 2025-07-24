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
		return []string{"id", "email", "contact_email", "first_name", "last_name", "phone", "address", "avatar", "role", "created_at", "updated_at"}
	}

	// Admin can see all fields for Professeur and Etudiant
	if viewerRole == RoleAdmin && (targetRole == RoleProfesseur || targetRole == RoleEtudiant) {
		return []string{"id", "email", "contact_email", "first_name", "last_name", "phone", "address", "avatar", "role", "created_at", "updated_at"}
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
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	ContactEmail string         `json:"contact_email"`
	Password     string         `json:"-" gorm:"not null"`
	FirstName    string         `json:"first_name" gorm:"not null"`
	LastName     string         `json:"last_name" gorm:"not null"`
	Phone        string         `json:"phone"`
	Address      string         `json:"address"`
	Avatar       string         `json:"avatar" gorm:"default:'/assets/images/avatars/default-avatar.png'"`
	Role         string         `json:"role" gorm:"default:'etudiant'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	ContactEmail string    `json:"contact_email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" binding:"email"`
	ContactEmail string `json:"contact_email" binding:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
}

type UpdateProfileRequest struct {
	ContactEmail string `json:"contact_email" binding:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// PasswordStrength represents the strength of a password
type PasswordStrength struct {
	Score    int    `json:"score"`    // 0-4 (0=très faible, 4=très fort)
	Feedback string `json:"feedback"` // Message d'aide
	IsValid  bool   `json:"is_valid"` // Si le mot de passe respecte les critères minimum
	Criteria struct {
		Length    bool `json:"length"`    // Au moins 8 caractères
		Uppercase bool `json:"uppercase"` // Au moins 1 majuscule
		Lowercase bool `json:"lowercase"` // Au moins 1 minuscule
		Number    bool `json:"number"`    // Au moins 1 chiffre
		Special   bool `json:"special"`   // Au moins 1 caractère spécial
	} `json:"criteria"`
}

// ValidatePasswordStrength validates password strength and returns detailed feedback
func ValidatePasswordStrength(password string) PasswordStrength {
	var strength PasswordStrength

	// Critères de base
	strength.Criteria.Length = len(password) >= 8
	strength.Criteria.Uppercase = hasUppercase(password)
	strength.Criteria.Lowercase = hasLowercase(password)
	strength.Criteria.Number = hasNumber(password)
	strength.Criteria.Special = hasSpecial(password)

	// Calcul du score
	score := 0
	if strength.Criteria.Length {
		score++
	}
	if strength.Criteria.Uppercase {
		score++
	}
	if strength.Criteria.Lowercase {
		score++
	}
	if strength.Criteria.Number {
		score++
	}
	if strength.Criteria.Special {
		score++
	}

	strength.Score = score
	strength.IsValid = score >= 4 // Au moins 4 critères sur 5

	// Messages de feedback
	switch score {
	case 0:
		strength.Feedback = "Mot de passe très faible"
	case 1:
		strength.Feedback = "Mot de passe faible"
	case 2:
		strength.Feedback = "Mot de passe moyen"
	case 3:
		strength.Feedback = "Mot de passe bon"
	case 4:
		strength.Feedback = "Mot de passe fort"
	case 5:
		strength.Feedback = "Mot de passe très fort"
	}

	return strength
}

// Helper functions for password validation
func hasUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func hasSpecial(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		for _, special := range specialChars {
			if r == special {
				return true
			}
		}
	}
	return false
}
