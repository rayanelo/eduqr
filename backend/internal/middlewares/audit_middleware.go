package middlewares

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuditMiddleware struct {
	auditLogService *services.AuditLogService
}

func NewAuditMiddleware(auditLogService *services.AuditLogService) *AuditMiddleware {
	return &AuditMiddleware{
		auditLogService: auditLogService,
	}
}

// AuditMiddleware logs user actions automatically
func (m *AuditMiddleware) AuditMiddleware(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user information from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userEmail, exists := c.Get("user_email")
		if !exists {
			c.Next()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.Next()
			return
		}

		// Get IP address
		ipAddress := c.ClientIP()
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Real-IP")
		}
		if ipAddress == "" {
			ipAddress = "unknown"
		}

		// Get User-Agent
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			userAgent = "unknown"
		}

		// Get resource ID from URL parameters
		var resourceID *uint
		if resourceType != "" {
			if idStr := c.Param("id"); idStr != "" {
				if id, err := parseUint(idStr); err == nil {
					resourceID = &id
				}
			}
		}

		// Create description based on action and resource type
		description := m.createDescription(action, resourceType, c)

		// Log the action
		go func() {
			_ = m.auditLogService.LogUserAction(
				userID.(uint),
				userEmail.(string),
				userRole.(string),
				action,
				resourceType,
				resourceID,
				description,
				ipAddress,
				userAgent,
				nil, // Old values - would need to be captured before the action
				nil, // New values - would need to be captured after the action
			)
		}()

		c.Next()
	}
}

// AuditLoginMiddleware logs user login
func (m *AuditMiddleware) AuditLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This will be called after successful login
		// User information will be available in the response
		c.Next()

		// Check if login was successful (status 200/201)
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// Try to get user info from response
			if userData, exists := c.Get("login_user_data"); exists {
				if user, ok := userData.(map[string]interface{}); ok {
					if userID, ok := user["id"].(float64); ok {
						if userEmail, ok := user["email"].(string); ok {
							if userRole, ok := user["role"].(string); ok {
								// Get IP address
								ipAddress := c.ClientIP()
								if ipAddress == "" {
									ipAddress = c.GetHeader("X-Forwarded-For")
								}
								if ipAddress == "" {
									ipAddress = c.GetHeader("X-Real-IP")
								}
								if ipAddress == "" {
									ipAddress = "unknown"
								}

								// Get User-Agent
								userAgent := c.GetHeader("User-Agent")
								if userAgent == "" {
									userAgent = "unknown"
								}

								// Log login
								go func() {
									_ = m.auditLogService.LogLogin(
										uint(userID),
										userEmail,
										userRole,
										ipAddress,
										userAgent,
									)
								}()
							}
						}
					}
				}
			}
		}
	}
}

// AuditLogoutMiddleware logs user logout
func (m *AuditMiddleware) AuditLogoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user information from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userEmail, exists := c.Get("user_email")
		if !exists {
			c.Next()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.Next()
			return
		}

		// Get IP address
		ipAddress := c.ClientIP()
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Real-IP")
		}
		if ipAddress == "" {
			ipAddress = "unknown"
		}

		// Get User-Agent
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			userAgent = "unknown"
		}

		// Log logout
		go func() {
			_ = m.auditLogService.LogLogout(
				userID.(uint),
				userEmail.(string),
				userRole.(string),
				ipAddress,
				userAgent,
			)
		}()

		c.Next()
	}
}

// createDescription creates a human-readable description for the audit log
func (m *AuditMiddleware) createDescription(action, resourceType string, c *gin.Context) string {
	switch action {
	case models.ActionCreate:
		switch resourceType {
		case models.ResourceUser:
			return "Création d'un nouvel utilisateur"
		case models.ResourceRoom:
			return "Création d'une nouvelle salle"
		case models.ResourceSubject:
			return "Création d'une nouvelle matière"
		case models.ResourceCourse:
			return "Création d'un nouveau cours"
		case models.ResourceEvent:
			return "Création d'un nouvel événement"
		default:
			return "Création d'une ressource"
		}
	case models.ActionUpdate:
		switch resourceType {
		case models.ResourceUser:
			return "Modification d'un utilisateur"
		case models.ResourceRoom:
			return "Modification d'une salle"
		case models.ResourceSubject:
			return "Modification d'une matière"
		case models.ResourceCourse:
			return "Modification d'un cours"
		case models.ResourceEvent:
			return "Modification d'un événement"
		default:
			return "Modification d'une ressource"
		}
	case models.ActionDelete:
		switch resourceType {
		case models.ResourceUser:
			return "Suppression d'un utilisateur"
		case models.ResourceRoom:
			return "Suppression d'une salle"
		case models.ResourceSubject:
			return "Suppression d'une matière"
		case models.ResourceCourse:
			return "Suppression d'un cours"
		case models.ResourceEvent:
			return "Suppression d'un événement"
		default:
			return "Suppression d'une ressource"
		}
	default:
		return action + " sur " + resourceType
	}
}

// parseUint safely parses a string to uint
func parseUint(s string) (uint, error) {
	// Remove any leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Simple implementation - in production you might want to use strconv.ParseUint
	var result uint
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + uint(char-'0')
		} else {
			return 0, &strconv.NumError{Func: "ParseUint", Num: s, Err: strconv.ErrSyntax}
		}
	}
	return result, nil
}
