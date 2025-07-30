package repositories

import (
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		db: database.GetDB(),
	}
}

func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("User").First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) FindByUserID(userID uint) ([]models.Event, error) {
	var events []models.Event
	err := r.db.Preload("User").Where("user_id = ?", userID).Find(&events).Error
	return events, err
}

func (r *EventRepository) FindAll() ([]models.Event, error) {
	var events []models.Event
	err := r.db.Preload("User").Find(&events).Error
	return events, err
}

func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

func (r *EventRepository) FindByDateRange(startDate, endDate time.Time, userID uint) ([]models.Event, error) {
	var events []models.Event
	err := r.db.Preload("User").
		Where("user_id = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))",
			userID, startDate, endDate, startDate, endDate, startDate, endDate).
		Find(&events).Error
	return events, err
}
