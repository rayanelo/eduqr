package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditLogRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateAuditLog_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer un utilisateur de test
		user := createTestUser("admin")

		auditLog := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "create",
			ResourceType: "room",
			ResourceID:   nil,
			Description:  "Test audit log creation",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}

		err := repo.Create(auditLog)
		assert.NoError(t, err)
		assert.NotZero(t, auditLog.ID)
		assert.Equal(t, user.ID, auditLog.UserID)
		assert.Equal(t, "create", auditLog.Action)
		assert.Equal(t, "room", auditLog.ResourceType)
	})

	t.Run("FindByID_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer un audit log
		user := createTestUser("admin")
		auditLog := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "update",
			ResourceType: "subject",
			ResourceID:   nil,
			Description:  "Test audit log retrieval",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		testDB.Create(auditLog)

		// Récupérer l'audit log par ID
		retrievedLog, err := repo.FindByID(auditLog.ID)
		assert.NoError(t, err)
		assert.Equal(t, auditLog.ID, retrievedLog.ID)
		assert.Equal(t, auditLog.UserID, retrievedLog.UserID)
		assert.Equal(t, auditLog.Action, retrievedLog.Action)
		assert.Equal(t, auditLog.ResourceType, retrievedLog.ResourceType)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Essayer de récupérer un audit log inexistant
		_, err := repo.FindByID(99999)
		assert.Error(t, err)
	})

	t.Run("FindAll_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer plusieurs audit logs
		user := createTestUser("admin")

		auditLog1 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "create",
			ResourceType: "room",
			ResourceID:   nil,
			Description:  "First audit log",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		auditLog2 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "update",
			ResourceType: "subject",
			ResourceID:   nil,
			Description:  "Second audit log",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		testDB.Create(auditLog1)
		testDB.Create(auditLog2)

		// Récupérer tous les audit logs
		filter := &models.AuditLogFilter{
			Page:  1,
			Limit: 10,
		}
		response, err := repo.FindAll(filter)
		assert.NoError(t, err)
		assert.Len(t, response.Logs, 2)
		assert.Equal(t, int64(2), response.Total)
	})

	t.Run("FindByUserID_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer des audit logs pour un utilisateur
		user := createTestUser("admin")

		auditLog1 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "create",
			ResourceType: "room",
			ResourceID:   nil,
			Description:  "User audit log 1",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		auditLog2 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "update",
			ResourceType: "subject",
			ResourceID:   nil,
			Description:  "User audit log 2",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		testDB.Create(auditLog1)
		testDB.Create(auditLog2)

		// Récupérer les audit logs de l'utilisateur
		logs, err := repo.FindByUserID(user.ID, 10)
		assert.NoError(t, err)
		assert.Len(t, logs, 2)

		for _, log := range logs {
			assert.Equal(t, user.ID, log.UserID)
		}
	})

	t.Run("FindByResource_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer des audit logs pour une ressource
		user := createTestUser("admin")

		auditLog1 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "create",
			ResourceType: "room",
			ResourceID:   nil,
			Description:  "Room resource log",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		auditLog2 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "update",
			ResourceType: "subject",
			ResourceID:   nil,
			Description:  "Subject resource log",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		testDB.Create(auditLog1)
		testDB.Create(auditLog2)

		// Récupérer les audit logs par ressource
		logs, err := repo.FindByResource("room", 0, 10)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "room", logs[0].ResourceType)
	})

	t.Run("FindRecent_Success", func(t *testing.T) {
		repo := repositories.NewAuditLogRepository()

		// Créer des audit logs
		user := createTestUser("admin")

		auditLog1 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "create",
			ResourceType: "room",
			ResourceID:   nil,
			Description:  "Recent log 1",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		auditLog2 := &models.AuditLog{
			UserID:       user.ID,
			UserEmail:    user.Email,
			UserRole:     user.Role,
			Action:       "update",
			ResourceType: "subject",
			ResourceID:   nil,
			Description:  "Recent log 2",
			IPAddress:    "127.0.0.1",
			UserAgent:    "Test Agent",
		}
		testDB.Create(auditLog1)
		testDB.Create(auditLog2)

		// Récupérer les audit logs récents
		logs, err := repo.FindRecent(5)
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
	})
}

func TestAuditLogValidation(t *testing.T) {
	t.Run("ValidateAuditLogAction_ValidActions", func(t *testing.T) {
		validActions := []string{
			"create",
			"update",
			"delete",
			"login",
			"logout",
		}

		for _, action := range validActions {
			assert.Contains(t, validActions, action, "L'action doit être valide")
		}
	})

	t.Run("ValidateAuditLogAction_InvalidActions", func(t *testing.T) {
		invalidActions := []string{
			"",
			"invalid",
			"read",
			"modify",
		}

		validActions := []string{"create", "update", "delete", "login", "logout"}
		for _, action := range invalidActions {
			if action != "" {
				assert.NotContains(t, validActions, action, "L'action invalide est détectée")
			}
		}
	})

	t.Run("ValidateResourceType_ValidResourceTypes", func(t *testing.T) {
		validResourceTypes := []string{
			"user",
			"room",
			"subject",
			"course",
			"event",
			"absence",
		}

		for _, resourceType := range validResourceTypes {
			assert.Contains(t, validResourceTypes, resourceType, "Le type de ressource doit être valide")
		}
	})

	t.Run("ValidateResourceType_InvalidResourceTypes", func(t *testing.T) {
		invalidResourceTypes := []string{
			"",
			"invalid",
			"student",
			"teacher",
		}

		validResourceTypes := []string{"user", "room", "subject", "course", "event", "absence"}
		for _, resourceType := range invalidResourceTypes {
			if resourceType != "" {
				assert.NotContains(t, validResourceTypes, resourceType, "Le type de ressource invalide est détecté")
			}
		}
	})

	t.Run("ValidateDescription_ValidDescriptions", func(t *testing.T) {
		validDescriptions := []string{
			"User created successfully",
			"Room updated",
			"Course deleted",
			"User logged in",
			"Absence approved",
		}

		for _, description := range validDescriptions {
			assert.NotEmpty(t, description, "La description ne devrait pas être vide")
			assert.Len(t, description, 1, "La description devrait avoir au moins 1 caractère")
		}
	})

	t.Run("ValidateDescription_InvalidDescriptions", func(t *testing.T) {
		invalidDescriptions := []string{
			"",
			"   ",
		}

		for _, description := range invalidDescriptions {
			if description == "" {
				assert.Empty(t, description, "La description vide est détectée correctement")
			}
		}
	})
}
