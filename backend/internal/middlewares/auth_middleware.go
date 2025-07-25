package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// AuthMiddleware validates JWT token and sets user information in context
func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := utils.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role or higher
func (m *AuthMiddleware) RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has the required role or higher
		if !hasRoleOrHigher(userRole.(string), requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CanManageUserMiddleware checks if the current user can manage the target user
func (m *AuthMiddleware) CanManageUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get target user ID from URL parameter
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID required"})
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		// Get current user ID
		currentUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// If user is trying to manage themselves, allow it
		if currentUserID.(uint) == uint(id) {
			c.Next()
			return
		}

		// Set target user ID and user role in context for controller to check permissions
		c.Set("target_user_id", uint(id))
		c.Set("current_user_role", userRole.(string))
		c.Next()
	}
}

// CanViewUserMiddleware checks if the current user can view the target user
func (m *AuthMiddleware) CanViewUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get target user ID from URL parameter
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID required"})
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		// Get current user ID
		currentUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// If user is trying to view themselves, allow it
		if currentUserID.(uint) == uint(id) {
			c.Next()
			return
		}

		// Set target user ID and user role in context for controller to check permissions
		c.Set("target_user_id", uint(id))
		c.Set("current_user_role", userRole.(string))
		c.Next()
	}
}

// hasRoleOrHigher checks if userRole has the same or higher privileges than requiredRole
func hasRoleOrHigher(userRole, requiredRole string) bool {
	userLevel := models.RoleHierarchy[userRole]
	requiredLevel := models.RoleHierarchy[requiredRole]
	return userLevel >= requiredLevel
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't require authentication
func (m *AuthMiddleware) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			c.Next()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// CanDeleteMiddleware checks if the current user can delete the target resource
func (m *AuthMiddleware) CanDeleteMiddleware(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Seuls les Admins et Super Admins peuvent supprimer
		if userRole != models.RoleAdmin && userRole != models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions for deletion"})
			c.Abort()
			return
		}

		// Pour les suppressions d'utilisateurs, vérifier les règles spéciales
		if resourceType == "user" {
			// Get target user ID from URL parameter
			idStr := c.Param("id")
			if idStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user ID required"})
				c.Abort()
				return
			}

			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
				c.Abort()
				return
			}

			// Get current user ID
			currentUserID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				c.Abort()
				return
			}

			// Un utilisateur ne peut pas se supprimer lui-même
			if currentUserID.(uint) == uint(id) {
				c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete your own account"})
				c.Abort()
				return
			}

			// Set target user ID in context for controller to check specific permissions
			c.Set("target_user_id", uint(id))
		}

		// Set resource type in context
		c.Set("resource_type", resourceType)
		c.Set("current_user_role", userRole.(string))
		c.Next()
	}
}

// CanDeleteAdminMiddleware checks if the current user can delete an admin
func (m *AuthMiddleware) CanDeleteAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Seul le Super Admin peut supprimer un Admin
		if userRole != models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "only super admin can delete admin accounts"})
			c.Abort()
			return
		}

		c.Next()
	}
}
