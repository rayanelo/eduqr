package controllers

import (
	"net/http"
	"strconv"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuditLogController struct {
	auditLogService *services.AuditLogService
}

func NewAuditLogController(auditLogService *services.AuditLogService) *AuditLogController {
	return &AuditLogController{
		auditLogService: auditLogService,
	}
}

// GetAuditLogs retrieves audit logs with filtering and pagination
// @Summary Get audit logs
// @Description Retrieve audit logs with filtering and pagination
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param action query string false "Filter by action (create, update, delete, login, logout)"
// @Param resource_type query string false "Filter by resource type (user, room, subject, course, event)"
// @Param resource_id query int false "Filter by resource ID"
// @Param user_id query int false "Filter by user ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search in description, user_email, user_role"
// @Success 200 {object} models.AuditLogListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs [get]
func (c *AuditLogController) GetAuditLogs(ctx *gin.Context) {
	// Parse query parameters
	filter := &models.AuditLogFilter{}

	// Pagination
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		} else {
			filter.Page = 1
		}
	} else {
		filter.Page = 1
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		} else {
			filter.Limit = 20
		}
	} else {
		filter.Limit = 20
	}

	// Filters
	filter.Action = ctx.Query("action")
	filter.ResourceType = ctx.Query("resource_type")
	filter.Search = ctx.Query("search")

	// Resource ID filter
	if resourceIDStr := ctx.Query("resource_id"); resourceIDStr != "" {
		if resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32); err == nil {
			resourceIDUint := uint(resourceID)
			filter.ResourceID = &resourceIDUint
		}
	}

	// User ID filter
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			filter.UserID = &userIDUint
		}
	}

	// Date filters
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set end date to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.EndDate = &endDate
		}
	}

	// Get audit logs
	result, err := c.auditLogService.GetAuditLogs(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetAuditLogByID retrieves a specific audit log by ID
// @Summary Get audit log by ID
// @Description Retrieve a specific audit log entry by ID
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Audit log ID"
// @Success 200 {object} models.AuditLogResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/audit-logs/{id} [get]
func (c *AuditLogController) GetAuditLogByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	auditLog, err := c.auditLogService.GetAuditLogByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		return
	}

	ctx.JSON(http.StatusOK, auditLog)
}

// GetAuditLogStats retrieves audit log statistics
// @Summary Get audit log statistics
// @Description Retrieve audit log statistics for a date range
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD, default: 30 days ago)"
// @Param end_date query string false "End date (YYYY-MM-DD, default: today)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs/stats [get]
func (c *AuditLogController) GetAuditLogStats(ctx *gin.Context) {
	// Parse date parameters
	startDate := time.Now().AddDate(0, 0, -30) // Default to 30 days ago
	endDate := time.Now()

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set end date to end of day
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	stats, err := c.auditLogService.GetStats(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit log statistics"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetRecentAuditLogs retrieves recent audit logs
// @Summary Get recent audit logs
// @Description Retrieve recent audit logs (last 10 by default)
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of recent logs (default: 10, max: 100)"
// @Success 200 {array} models.AuditLogResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs/recent [get]
func (c *AuditLogController) GetRecentAuditLogs(ctx *gin.Context) {
	limit := 10 // Default limit

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	logs, err := c.auditLogService.GetRecentLogs(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recent audit logs"})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

// GetUserActivity retrieves activity for a specific user
// @Summary Get user activity
// @Description Retrieve activity summary for a specific user
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Param days query int false "Number of days (default: 30, max: 365)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs/user/{user_id}/activity [get]
func (c *AuditLogController) GetUserActivity(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	days := 30 // Default to 30 days
	if daysStr := ctx.Query("days"); daysStr != "" {
		if parsed, err := strconv.Atoi(daysStr); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	activity, err := c.auditLogService.GetUserActivity(uint(userID), days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user activity"})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// GetResourceHistory retrieves audit logs for a specific resource
// @Summary Get resource history
// @Description Retrieve audit logs for a specific resource
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_type path string true "Resource type (user, room, subject, course, event)"
// @Param resource_id path int true "Resource ID"
// @Param limit query int false "Number of logs (default: 20, max: 100)"
// @Success 200 {array} models.AuditLogResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs/resource/{resource_type}/{resource_id} [get]
func (c *AuditLogController) GetResourceHistory(ctx *gin.Context) {
	resourceType := ctx.Param("resource_type")
	resourceIDStr := ctx.Param("resource_id")

	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	limit := 20 // Default limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	logs, err := c.auditLogService.GetResourceHistory(resourceType, uint(resourceID), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve resource history"})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

// CleanOldLogs removes audit logs older than the specified duration
// @Summary Clean old audit logs
// @Description Remove audit logs older than the specified duration
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int true "Remove logs older than N days"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/audit-logs/clean [delete]
func (c *AuditLogController) CleanOldLogs(ctx *gin.Context) {
	daysStr := ctx.Query("days")
	if daysStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Days parameter is required"})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
		return
	}

	// Limit to max 365 days to prevent accidental deletion of all logs
	if days > 365 {
		days = 365
	}

	olderThan := time.Now().AddDate(0, 0, -days)
	err = c.auditLogService.CleanOldLogs(olderThan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean old audit logs"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Old audit logs cleaned successfully",
		"older_than": olderThan.Format("2006-01-02 15:04:05"),
		"days":       days,
	})
}
