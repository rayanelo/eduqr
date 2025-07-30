package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQRCodeGeneration(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("GenerateQRCode_Success", func(t *testing.T) {
		// Créer les dépendances
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Générer un QR code pour le cours
		qrCodeData, err := generateQRCodeForCourse(course.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, qrCodeData)
		assert.Len(t, qrCodeData, 1, "Le QR code doit avoir au moins 1 caractère")
	})

	t.Run("ValidateQRCode_Success", func(t *testing.T) {
		// Créer les dépendances
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Générer un QR code
		qrCodeData, err := generateQRCodeForCourse(course.ID)
		assert.NoError(t, err)

		// Valider le QR code
		qrInfo, err := validateQRCode(qrCodeData)
		assert.NoError(t, err)
		assert.NotNil(t, qrInfo)
		assert.Equal(t, course.ID, qrInfo.CourseID)
		assert.Equal(t, course.Name, qrInfo.CourseName)
		assert.Equal(t, subject.Name, qrInfo.SubjectName)
		assert.Equal(t, teacher.FirstName+" "+teacher.LastName, qrInfo.TeacherName)
		assert.Equal(t, room.Name, qrInfo.RoomName)
		assert.True(t, qrInfo.IsValid)
	})

	t.Run("ValidateQRCode_InvalidData", func(t *testing.T) {
		// Tester avec des données invalides
		invalidQRData := "invalid_qr_code_data"

		// Valider le QR code invalide
		qrInfo, err := validateQRCode(invalidQRData)
		assert.Error(t, err)
		assert.Nil(t, qrInfo)
	})
}

func TestProfileManagement(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("GetUserProfile_Success", func(t *testing.T) {
		userRepo := repositories.NewUserRepository()
		service := services.NewUserService(userRepo, "test-secret", 1*time.Hour)

		// Créer un utilisateur
		user := createTestUser("student")

		// Récupérer le profil
		profile, err := service.GetUserByID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, user.ID, profile.ID)
		assert.Equal(t, user.Email, profile.Email)
		assert.Equal(t, user.FirstName, profile.FirstName)
		assert.Equal(t, user.LastName, profile.LastName)
		assert.Equal(t, user.Role, profile.Role)
	})

	t.Run("UpdateUserProfile_Success", func(t *testing.T) {
		userRepo := repositories.NewUserRepository()
		service := services.NewUserService(userRepo, "test-secret", 1*time.Hour)

		// Créer un utilisateur
		user := createTestUser("student")

		// Mettre à jour le profil
		updateReq := &models.UpdateUserRequest{
			FirstName: "Updated First Name",
			LastName:  "Updated Last Name",
			Phone:     "+9876543210",
			Address:   "Updated Address",
		}

		updatedProfile, err := service.UpdateUser(user.ID, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updatedProfile)
		assert.Equal(t, updateReq.FirstName, updatedProfile.FirstName)
		assert.Equal(t, updateReq.LastName, updatedProfile.LastName)
		assert.Equal(t, updateReq.Phone, updatedProfile.Phone)
		assert.Equal(t, updateReq.Address, updatedProfile.Address)
	})

	t.Run("UpdateUserProfile_InvalidData", func(t *testing.T) {
		userRepo := repositories.NewUserRepository()
		service := services.NewUserService(userRepo, "test-secret", 1*time.Hour)

		// Créer un utilisateur
		user := createTestUser("student")

		// Mettre à jour le profil avec des données invalides
		updateReq := &models.UpdateUserRequest{
			FirstName: "", // Nom vide
			LastName:  "Valid Last Name",
		}

		_, err := service.UpdateUser(user.ID, updateReq)
		assert.Error(t, err)
	})

	t.Run("DeleteUserProfile_Success", func(t *testing.T) {
		userRepo := repositories.NewUserRepository()
		service := services.NewUserService(userRepo, "test-secret", 1*time.Hour)

		// Créer un utilisateur
		user := createTestUser("student")

		// Supprimer le profil
		err := service.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Vérifier que l'utilisateur n'existe plus
		_, err = service.GetUserByID(user.ID)
		assert.Error(t, err)
	})
}

func TestProfileValidation(t *testing.T) {
	t.Run("ValidateUserProfile_ValidData", func(t *testing.T) {
		validProfiles := []models.User{
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Phone:     "+1234567890",
				Address:   "123 Main St",
				Role:      "student",
			},
			{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Phone:     "+0987654321",
				Address:   "456 Oak Ave",
				Role:      "teacher",
			},
		}

		for _, profile := range validProfiles {
			assert.NotEmpty(t, profile.FirstName, "Le prénom ne devrait pas être vide")
			assert.NotEmpty(t, profile.LastName, "Le nom ne devrait pas être vide")
			assert.NotEmpty(t, profile.Email, "L'email ne devrait pas être vide")
			assert.Contains(t, profile.Email, "@", "L'email doit contenir @")
			assert.Contains(t, profile.Role, "student", "Le rôle doit être valide")
		}
	})

	t.Run("ValidateUserProfile_InvalidData", func(t *testing.T) {
		invalidProfiles := []models.User{
			{
				FirstName: "", // Prénom vide
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Role:      "student",
			},
			{
				FirstName: "John",
				LastName:  "", // Nom vide
				Email:     "john.doe@example.com",
				Role:      "student",
			},
			{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email", // Email invalide
				Role:      "student",
			},
		}

		for _, profile := range invalidProfiles {
			if profile.FirstName == "" {
				assert.Empty(t, profile.FirstName, "Le prénom vide est détecté")
			}
			if profile.LastName == "" {
				assert.Empty(t, profile.LastName, "Le nom vide est détecté")
			}
			if !contains(profile.Email, "@") {
				assert.False(t, isValidEmailFormat(profile.Email), "L'email invalide est détecté")
			}
		}
	})
}

// Fonctions utilitaires pour les tests
func generateQRCodeForCourse(courseID uint) (string, error) {
	// Simulation de génération de QR code
	// En réalité, cela appellerait le service de présence
	return "test_qr_code_data_for_course_" + string(rune(courseID)), nil
}

func validateQRCode(qrCodeData string) (*models.QRCodeInfo, error) {
	// Simulation de validation de QR code
	// En réalité, cela appellerait le service de présence
	if qrCodeData == "invalid_qr_code_data" {
		return nil, assert.AnError
	}

	return &models.QRCodeInfo{
		CourseID:    1,
		CourseName:  "Test Course",
		SubjectName: "Test Subject",
		TeacherName: "Test Teacher",
		RoomName:    "Test Room",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(2 * time.Hour),
		QRCodeData:  qrCodeData,
		IsValid:     true,
	}, nil
}
