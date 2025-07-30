package services

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"errors"
	"time"
)

type EventService struct {
	eventRepo *repositories.EventRepository
}

func NewEventService(eventRepo *repositories.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (s *EventService) CreateEvent(userID uint, req *models.CreateEventRequest) (*models.EventResponse, error) {
	// Validate time range
	if req.StartTime.After(req.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Color:       req.Color,
		UserID:      userID,
	}

	if event.Color == "" {
		event.Color = "#2196F3"
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	return s.toEventResponse(event), nil
}

func (s *EventService) GetEventByID(id uint, userID uint) (*models.EventResponse, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the event
	if event.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return s.toEventResponse(event), nil
}

func (s *EventService) GetUserEvents(userID uint) ([]models.EventResponse, error) {
	events, err := s.eventRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.EventResponse
	for _, event := range events {
		responses = append(responses, *s.toEventResponse(&event))
	}

	return responses, nil
}

func (s *EventService) UpdateEvent(id uint, userID uint, req *models.UpdateEventRequest) (*models.EventResponse, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the event
	if event.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Validate time range if both times are provided
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		if req.StartTime.After(req.EndTime) {
			return nil, errors.New("start time must be before end time")
		}
	}

	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if !req.StartTime.IsZero() {
		event.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		event.EndTime = req.EndTime
	}
	if req.Color != "" {
		event.Color = req.Color
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	return s.toEventResponse(event), nil
}

func (s *EventService) DeleteEvent(id uint, userID uint) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if user owns the event
	if event.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.eventRepo.Delete(id)
}

func (s *EventService) GetEventsByDateRange(userID uint, startDate, endDate time.Time) ([]models.EventResponse, error) {
	events, err := s.eventRepo.FindByDateRange(startDate, endDate, userID)
	if err != nil {
		return nil, err
	}

	var responses []models.EventResponse
	for _, event := range events {
		responses = append(responses, *s.toEventResponse(&event))
	}

	return responses, nil
}

func (s *EventService) toEventResponse(event *models.Event) *models.EventResponse {
	return &models.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Color:       event.Color,
		UserID:      event.UserID,
		User: models.UserResponse{
			ID:        event.User.ID,
			Email:     event.User.Email,
			FirstName: event.User.FirstName,
			LastName:  event.User.LastName,
			Role:      event.User.Role,
			CreatedAt: event.User.CreatedAt,
			UpdatedAt: event.User.UpdatedAt,
		},
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}
