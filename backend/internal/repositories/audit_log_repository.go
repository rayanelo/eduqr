package repositories

import (
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{
		db: database.GetDB(),
	}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(auditLog *models.AuditLog) error {
	return r.db.Create(auditLog).Error
}

// FindByID finds an audit log by ID
func (r *AuditLogRepository) FindByID(id uint) (*models.AuditLog, error) {
	var auditLog models.AuditLog
	err := r.db.First(&auditLog, id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// FindAll retrieves all audit logs with pagination and filtering
func (r *AuditLogRepository) FindAll(filter *models.AuditLogFilter) (*models.AuditLogListResponse, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}

	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}

	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(description) LIKE ? OR LOWER(user_email) LIKE ? OR LOWER(user_role) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (filter.Page - 1) * filter.Limit
	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	// Get paginated results
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responses []models.AuditLogResponse
	for _, log := range logs {
		responses = append(responses, r.toAuditLogResponse(&log))
	}

	return &models.AuditLogListResponse{
		Logs:       responses,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// FindByUserID retrieves audit logs for a specific user
func (r *AuditLogRepository) FindByUserID(userID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindByResource retrieves audit logs for a specific resource
func (r *AuditLogRepository) FindByResource(resourceType string, resourceID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindRecent retrieves recent audit logs
func (r *AuditLogRepository) FindRecent(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetStats retrieves audit log statistics
func (r *AuditLogRepository) GetStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total logs in date range
	var totalLogs int64
	err := r.db.Model(&models.AuditLog{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalLogs).Error
	if err != nil {
		return nil, err
	}
	stats["total_logs"] = totalLogs

	// Logs by action
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	err = r.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("action").
		Find(&actionStats).Error
	if err != nil {
		return nil, err
	}
	stats["by_action"] = actionStats

	// Logs by resource type
	var resourceStats []struct {
		ResourceType string `json:"resource_type"`
		Count        int64  `json:"count"`
	}
	err = r.db.Model(&models.AuditLog{}).
		Select("resource_type, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("resource_type").
		Find(&resourceStats).Error
	if err != nil {
		return nil, err
	}
	stats["by_resource"] = resourceStats

	// Logs by user role
	var roleStats []struct {
		UserRole string `json:"user_role"`
		Count    int64  `json:"count"`
	}
	err = r.db.Model(&models.AuditLog{}).
		Select("user_role, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("user_role").
		Find(&roleStats).Error
	if err != nil {
		return nil, err
	}
	stats["by_role"] = roleStats

	return stats, nil
}

// CleanOldLogs removes audit logs older than the specified duration
func (r *AuditLogRepository) CleanOldLogs(olderThan time.Time) error {
	return r.db.Where("created_at < ?", olderThan).Delete(&models.AuditLog{}).Error
}

// toAuditLogResponse converts AuditLog to AuditLogResponse
func (r *AuditLogRepository) toAuditLogResponse(log *models.AuditLog) models.AuditLogResponse {
	return models.AuditLogResponse{
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

// GetUserActivity retrieves user activity summary
func (r *AuditLogRepository) GetUserActivity(userID uint, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var activity []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	err := r.db.Model(&models.AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("DATE(created_at)").
		Order("date DESC").
		Find(&activity).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":  userID,
		"days":     days,
		"activity": activity,
	}, nil
}
