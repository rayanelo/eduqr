package repositories

import (
	"eduqr-backend/internal/models"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// GetAllRooms récupère toutes les salles avec filtres
func (r *RoomRepository) GetAllRooms(filter *models.RoomFilter) ([]models.Room, error) {
	var rooms []models.Room
	query := r.db.Preload("Parent").Preload("Children")

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.Building != "" {
			query = query.Where("building ILIKE ?", "%"+filter.Building+"%")
		}
		if filter.Floor != "" {
			query = query.Where("floor ILIKE ?", "%"+filter.Floor+"%")
		}
		if filter.Modular != nil {
			query = query.Where("is_modular = ?", *filter.Modular)
		}
		if filter.ParentID != nil {
			query = query.Where("parent_id = ?", *filter.ParentID)
		} else {
			// Par défaut, ne montrer que les salles principales (pas les sous-salles)
			query = query.Where("parent_id IS NULL")
		}
	}

	err := query.Find(&rooms).Error
	return rooms, err
}

// GetRoomByID récupère une salle par son ID
func (r *RoomRepository) GetRoomByID(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.Preload("Parent").Preload("Children").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// GetRoomByName récupère une salle par son nom
func (r *RoomRepository) GetRoomByName(name string) (*models.Room, error) {
	var room models.Room
	err := r.db.Where("name = ?", name).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// CreateRoom crée une nouvelle salle
func (r *RoomRepository) CreateRoom(room *models.Room) error {
	return r.db.Create(room).Error
}

// UpdateRoom met à jour une salle
func (r *RoomRepository) UpdateRoom(room *models.Room) error {
	return r.db.Save(room).Error
}

// DeleteRoom supprime une salle (soft delete)
func (r *RoomRepository) DeleteRoom(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// GetRoomsByParentID récupère toutes les sous-salles d'une salle parent
func (r *RoomRepository) GetRoomsByParentID(parentID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Where("parent_id = ?", parentID).Find(&rooms).Error
	return rooms, err
}

// CheckRoomExists vérifie si une salle existe déjà avec le même nom
func (r *RoomRepository) CheckRoomExists(name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Room{}).Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetModularRooms récupère toutes les salles modulables
func (r *RoomRepository) GetModularRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Where("is_modular = ? AND parent_id IS NULL", true).Find(&rooms).Error
	return rooms, err
}

// GetChildRooms récupère les salles enfants d'une salle modulable
func (r *RoomRepository) GetChildRooms(parentID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Where("parent_id = ?", parentID).Find(&rooms).Error
	return rooms, err
}
