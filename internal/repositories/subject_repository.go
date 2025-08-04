package repositories

import (
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"

	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository() *SubjectRepository {
	return &SubjectRepository{
		db: database.GetDB(),
	}
}

// GetAllSubjects récupère toutes les matières
func (r *SubjectRepository) GetAllSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Find(&subjects).Error
	return subjects, err
}

// GetSubjectByID récupère une matière par son ID
func (r *SubjectRepository) GetSubjectByID(id uint) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.First(&subject, id).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// GetSubjectByName récupère une matière par son nom
func (r *SubjectRepository) GetSubjectByName(name string) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Where("name = ?", name).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// GetSubjectByCode récupère une matière par son code
func (r *SubjectRepository) GetSubjectByCode(code string) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Where("code = ?", code).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// CreateSubject crée une nouvelle matière
func (r *SubjectRepository) CreateSubject(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

// UpdateSubject met à jour une matière
func (r *SubjectRepository) UpdateSubject(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

// DeleteSubject supprime une matière (soft delete)
func (r *SubjectRepository) DeleteSubject(id uint) error {
	return r.db.Delete(&models.Subject{}, id).Error
}

// CheckSubjectExists vérifie si une matière existe déjà avec le même nom
func (r *SubjectRepository) CheckSubjectExists(name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Subject{}).Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CheckSubjectInUse vérifie si une matière est utilisée dans des cours
func (r *SubjectRepository) CheckSubjectInUse(id uint) (bool, error) {
	// TODO: Implémenter quand la table des cours sera créée
	// Pour l'instant, on retourne false
	return false, nil
}

// CheckSubjectCodeExists vérifie si une matière existe déjà avec le même code
func (r *SubjectRepository) CheckSubjectCodeExists(code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Subject{}).Where("code = ?", code)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
