package services

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"errors"
	"time"
)

type AuditLogService struct {
	auditLogRepo *repositories.AuditLogRepository
}

func NewAuditLogService(auditLogRepo *repositories.AuditLogRepository) *AuditLogService {
	return &AuditLogService{
		auditLogRepo: auditLogRepo,
	}
}

// CreateAuditLog creates a new audit log entry
func (s *AuditLogService) CreateAuditLog(req *models.AuditLogRequest) (*models.AuditLogResponse, error) {
	// Validate action
	if !isValidAction(req.Action) {
		return nil, errors.New("invalid action")
	}

	// Validate resource type
	if !isValidResourceType(req.ResourceType) {
		return nil, errors.New("invalid resource type")
	}

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:       req.UserID,
		UserEmail:    req.UserEmail,
		UserRole:     req.UserRole,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Description:  req.Description,
		OldValues:    models.ToJSONString(req.OldValues),
		NewValues:    models.ToJSONString(req.NewValues),
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
	}

	if err := s.auditLogRepo.Create(auditLog); err != nil {
		return nil, err
	}

	return s.toAuditLogResponse(auditLog), nil
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (s *AuditLogService) GetAuditLogs(filter *models.AuditLogFilter) (*models.AuditLogListResponse, error) {
	// Set default values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	return s.auditLogRepo.FindAll(filter)
}

// GetAuditLogByID retrieves a specific audit log by ID
func (s *AuditLogService) GetAuditLogByID(id uint) (*models.AuditLogResponse, error) {
	auditLog, err := s.auditLogRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toAuditLogResponse(auditLog), nil
}

// GetUserActivity retrieves activity for a specific user
func (s *AuditLogService) GetUserActivity(userID uint, days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 30 // Default to 30 days
	}
	if days > 365 {
		days = 365 // Max 1 year
	}

	return s.auditLogRepo.GetUserActivity(userID, days)
}

// GetStats retrieves audit log statistics
func (s *AuditLogService) GetStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Limit to max 1 year
	if endDate.Sub(startDate) > 365*24*time.Hour {
		endDate = startDate.Add(365 * 24 * time.Hour)
	}

	return s.auditLogRepo.GetStats(startDate, endDate)
}

// GetRecentLogs retrieves recent audit logs
func (s *AuditLogService) GetRecentLogs(limit int) ([]models.AuditLogResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.auditLogRepo.FindRecent(limit)
	if err != nil {
		return nil, err
	}

	var responses []models.AuditLogResponse
	for _, log := range logs {
		responses = append(responses, *s.toAuditLogResponse(&log))
	}

	return responses, nil
}

// GetResourceHistory retrieves audit logs for a specific resource
func (s *AuditLogService) GetResourceHistory(resourceType string, resourceID uint, limit int) ([]models.AuditLogResponse, error) {
	if !isValidResourceType(resourceType) {
		return nil, errors.New("invalid resource type")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.auditLogRepo.FindByResource(resourceType, resourceID, limit)
	if err != nil {
		return nil, err
	}

	var responses []models.AuditLogResponse
	for _, log := range logs {
		responses = append(responses, *s.toAuditLogResponse(&log))
	}

	return responses, nil
}

// CleanOldLogs removes audit logs older than the specified duration
func (s *AuditLogService) CleanOldLogs(olderThan time.Time) error {
	return s.auditLogRepo.CleanOldLogs(olderThan)
}

// LogUserAction is a convenience method to log user actions
func (s *AuditLogService) LogUserAction(userID uint, userEmail, userRole, action, resourceType string, resourceID *uint, description, ipAddress, userAgent string, oldValues, newValues interface{}) error {
	req := &models.AuditLogRequest{
		UserID:       userID,
		UserEmail:    userEmail,
		UserRole:     userRole,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
		OldValues:    oldValues,
		NewValues:    newValues,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	_, err := s.CreateAuditLog(req)
	return err
}

// LogLogin logs a user login
func (s *AuditLogService) LogLogin(userID uint, userEmail, userRole, ipAddress, userAgent string) error {
	return s.LogUserAction(
		userID,
		userEmail,
		userRole,
		models.ActionLogin,
		"",  // No resource type for login
		nil, // No resource ID for login
		"Connexion utilisateur",
		ipAddress,
		userAgent,
		nil,
		nil,
	)
}

// LogLogout logs a user logout
func (s *AuditLogService) LogLogout(userID uint, userEmail, userRole, ipAddress, userAgent string) error {
	return s.LogUserAction(
		userID,
		userEmail,
		userRole,
		models.ActionLogout,
		"",  // No resource type for logout
		nil, // No resource ID for logout
		"Déconnexion utilisateur",
		ipAddress,
		userAgent,
		nil,
		nil,
	)
}

// toAuditLogResponse converts AuditLog to AuditLogResponse
func (s *AuditLogService) toAuditLogResponse(log *models.AuditLog) *models.AuditLogResponse {
	return &models.AuditLogResponse{
		ID:           log.ID,
		UserID:       log.UserID,
		UserEmail:    log.UserEmail,
		UserRole:     log.UserRole,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Description:  log.Description,
		OldValues:    log.OldValues,
		NewValues:    log.NewValues,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		CreatedAt:    log.CreatedAt,
	}
}

// Helper functions for validation
func isValidAction(action string) bool {
	validActions := []string{
		models.ActionCreate,
		models.ActionUpdate,
		models.ActionDelete,
		models.ActionLogin,
		models.ActionLogout,
	}

	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}

func isValidResourceType(resourceType string) bool {
	validResourceTypes := []string{
		models.ResourceUser,
		models.ResourceRoom,
		models.ResourceSubject,
		models.ResourceCourse,
		models.ResourceEvent,
	}

	for _, validType := range validResourceTypes {
		if resourceType == validType {
			return true
		}
	}
	return false
}
