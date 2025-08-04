package repositories

import (
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// DeleteUser supprime un utilisateur (soft delete)
func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// CreateUser alias pour Create
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.Create(user)
}

// GetUserByID alias pour FindByID
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	return r.FindByID(id)
}

// GetUserByEmail alias pour FindByEmail
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	return r.FindByEmail(email)
}
