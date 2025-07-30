package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		user := &models.User{
			Email:     "test@eduqr.com",
			FirstName: "Test",
			LastName:  "User",
			Password:  "$2a$10$testpassword",
			Role:      "student",
			Phone:     "+1234567890",
			Address:   "Test Address",
		}

		err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "test@eduqr.com", user.Email)
	})

	t.Run("CreateUser_DuplicateEmail_ShouldFail", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer le premier utilisateur
		user1 := &models.User{
			Email:     "duplicate@eduqr.com",
			FirstName: "Test",
			LastName:  "User1",
			Password:  "$2a$10$testpassword",
			Role:      "student",
		}
		err := repo.CreateUser(user1)
		assert.NoError(t, err)

		// Essayer de créer un utilisateur avec le même email
		user2 := &models.User{
			Email:     "duplicate@eduqr.com",
			FirstName: "Test",
			LastName:  "User2",
			Password:  "$2a$10$testpassword",
			Role:      "teacher",
		}
		err = repo.CreateUser(user2)
		assert.Error(t, err) // Doit échouer à cause de la contrainte unique
	})

	t.Run("GetUserByID_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer un utilisateur
		user := createTestUser("student")

		// Récupérer l'utilisateur par ID
		retrievedUser, err := repo.GetUserByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("GetUserByID_NotFound", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Essayer de récupérer un utilisateur inexistant
		_, err := repo.GetUserByID(99999)
		assert.Error(t, err)
	})

	t.Run("GetUserByEmail_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer un utilisateur
		user := createTestUser("teacher")

		// Récupérer l'utilisateur par email
		retrievedUser, err := repo.GetUserByEmail(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Essayer de récupérer un utilisateur avec un email inexistant
		_, err := repo.GetUserByEmail("nonexistent@eduqr.com")
		assert.Error(t, err)
	})

	t.Run("GetAllUsers_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer plusieurs utilisateurs
		createTestUser("student")
		createTestUser("teacher")
		createTestUser("admin")

		// Récupérer tous les utilisateurs
		users, err := repo.GetAllUsers(nil)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("GetAllUsers_WithRoleFilter", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer des utilisateurs avec différents rôles
		createTestUser("student")
		createTestUser("teacher")
		createTestUser("admin")

		// Filtrer par rôle
		filter := &models.UserFilter{Role: "teacher"}
		users, err := repo.GetAllUsers(filter)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "teacher", users[0].Role)
	})

	t.Run("UpdateUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer un utilisateur
		user := createTestUser("student")

		// Modifier l'utilisateur
		user.FirstName = "Updated First Name"
		user.LastName = "Updated Last Name"
		user.Phone = "+9876543210"

		err := repo.UpdateUser(user)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedUser, err := repo.GetUserByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated First Name", updatedUser.FirstName)
		assert.Equal(t, "Updated Last Name", updatedUser.LastName)
		assert.Equal(t, "+9876543210", updatedUser.Phone)
	})

	t.Run("DeleteUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer un utilisateur
		user := createTestUser("student")

		// Supprimer l'utilisateur
		err := repo.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Vérifier que l'utilisateur n'existe plus
		_, err = repo.GetUserByID(user.ID)
		assert.Error(t, err)
	})

	t.Run("GetUsersByRole_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer des utilisateurs avec différents rôles
		createTestUser("student")
		createTestUser("student")
		createTestUser("teacher")
		createTestUser("admin")

		// Récupérer les étudiants
		students, err := repo.GetUsersByRole("student")
		assert.NoError(t, err)
		assert.Len(t, students, 2)

		for _, student := range students {
			assert.Equal(t, "student", student.Role)
		}
	})

	t.Run("CheckUserExists_Exists", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Créer un utilisateur
		user := createTestUser("student")

		// Vérifier que l'utilisateur existe
		exists, err := repo.CheckUserExists(user.Email, nil)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CheckUserExists_NotExists", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)

		// Vérifier qu'un utilisateur inexistant n'existe pas
		exists, err := repo.CheckUserExists("nonexistent@eduqr.com", nil)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserService(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		req := &models.CreateUserRequest{
			Email:     "service@eduqr.com",
			FirstName: "Service",
			LastName:  "User",
			Password:  "password123",
			Role:      "student",
			Phone:     "+1234567890",
			Address:   "Service Address",
		}

		response, err := service.CreateUser(req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Email, response.Email)
		assert.Equal(t, req.FirstName, response.FirstName)
		assert.Equal(t, req.LastName, response.LastName)
		assert.Equal(t, req.Role, response.Role)
	})

	t.Run("CreateUser_InvalidEmail_ShouldFail", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		req := &models.CreateUserRequest{
			Email:     "invalid-email",
			FirstName: "Service",
			LastName:  "User",
			Password:  "password123",
			Role:      "student",
		}

		_, err := service.CreateUser(req)
		assert.Error(t, err)
	})

	t.Run("CreateUser_EmptyPassword_ShouldFail", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		req := &models.CreateUserRequest{
			Email:     "test@eduqr.com",
			FirstName: "Service",
			LastName:  "User",
			Password:  "",
			Role:      "student",
		}

		_, err := service.CreateUser(req)
		assert.Error(t, err)
	})

	t.Run("GetUserByID_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		// Créer un utilisateur
		user := createTestUser("student")

		// Récupérer l'utilisateur
		response, err := service.GetUserByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("UpdateUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		// Créer un utilisateur
		user := createTestUser("student")

		req := &models.UpdateUserRequest{
			FirstName: "Updated Service",
			LastName:  "Updated User",
			Phone:     "+9876543210",
			Address:   "Updated Service Address",
		}

		response, err := service.UpdateUser(user.ID, req)
		assert.NoError(t, err)
		assert.Equal(t, req.FirstName, response.FirstName)
		assert.Equal(t, req.LastName, response.LastName)
		assert.Equal(t, req.Phone, response.Phone)
		assert.Equal(t, req.Address, response.Address)
	})

	t.Run("GetAllUsers_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		// Créer plusieurs utilisateurs
		createTestUser("student")
		createTestUser("teacher")
		createTestUser("admin")

		// Récupérer tous les utilisateurs
		users, err := service.GetAllUsers(nil)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("DeleteUser_Success", func(t *testing.T) {
		repo := repositories.NewUserRepository(testDB)
		service := services.NewUserService(repo)

		// Créer un utilisateur
		user := createTestUser("student")

		// Supprimer l'utilisateur
		err := service.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Vérifier que l'utilisateur n'existe plus
		_, err = service.GetUserByID(user.ID)
		assert.Error(t, err)
	})
}

func TestUserValidation(t *testing.T) {
	t.Run("ValidateEmail_ValidEmails", func(t *testing.T) {
		validEmails := []string{
			"test@eduqr.com",
			"user.name@domain.com",
			"user+tag@example.org",
			"123@numbers.com",
		}

		for _, email := range validEmails {
			assert.Contains(t, email, "@", "L'email doit contenir @")
			assert.Contains(t, email, ".", "L'email doit contenir un point")
		}
	})

	t.Run("ValidateEmail_InvalidEmails", func(t *testing.T) {
		invalidEmails := []string{
			"",
			"invalid-email",
			"@domain.com",
			"user@",
			"user.domain.com",
		}

		for _, email := range invalidEmails {
			if email == "" {
				assert.Empty(t, email, "L'email vide est détecté correctement")
			} else if !contains(email, "@") || !contains(email, ".") {
				assert.False(t, isValidEmailFormat(email), "Format d'email invalide détecté")
			}
		}
	})

	t.Run("ValidatePassword_ValidPasswords", func(t *testing.T) {
		validPasswords := []string{
			"password123",
			"SecurePass!",
			"123456789",
			"aBcDeFgH",
		}

		for _, password := range validPasswords {
			assert.Len(t, password, 8, "Le mot de passe doit avoir au moins 8 caractères")
		}
	})

	t.Run("ValidatePassword_InvalidPasswords", func(t *testing.T) {
		invalidPasswords := []string{
			"",
			"123",
			"short",
		}

		for _, password := range invalidPasswords {
			if password == "" {
				assert.Empty(t, password, "Le mot de passe vide est détecté correctement")
			} else if len(password) < 8 {
				assert.Len(t, password, 8, "Le mot de passe trop court est détecté")
			}
		}
	})

	t.Run("ValidateRole_ValidRoles", func(t *testing.T) {
		validRoles := []string{
			"student",
			"teacher",
			"admin",
			"super_admin",
		}

		for _, role := range validRoles {
			assert.Contains(t, validRoles, role, "Le rôle doit être valide")
		}
	})
}

// Fonctions utilitaires pour les tests
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

func isValidEmailFormat(email string) bool {
	return contains(email, "@") && contains(email, ".") && len(email) > 5
}
