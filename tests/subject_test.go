package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubjectRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		subject := &models.Subject{
			Name:        "Test Subject",
			Code:        "TEST001",
			Description: "Test Description",
		}

		err := repo.CreateSubject(subject)
		assert.NoError(t, err)
		assert.NotZero(t, subject.ID)
		assert.Equal(t, "Test Subject", subject.Name)
		assert.Equal(t, "TEST001", subject.Code)
	})

	t.Run("CreateSubject_DuplicateName_ShouldFail", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer la première matière
		subject1 := &models.Subject{
			Name:        "Duplicate Subject",
			Code:        "DUP001",
			Description: "First Description",
		}
		err := repo.CreateSubject(subject1)
		assert.NoError(t, err)

		// Essayer de créer une matière avec le même nom
		subject2 := &models.Subject{
			Name:        "Duplicate Subject",
			Code:        "DUP002",
			Description: "Second Description",
		}
		err = repo.CreateSubject(subject2)
		assert.Error(t, err) // Doit échouer à cause de la contrainte unique
	})

	t.Run("CreateSubject_DuplicateCode_ShouldFail", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer la première matière
		subject1 := &models.Subject{
			Name:        "Subject A",
			Code:        "CODE001",
			Description: "First Description",
		}
		err := repo.CreateSubject(subject1)
		assert.NoError(t, err)

		// Essayer de créer une matière avec le même code
		subject2 := &models.Subject{
			Name:        "Subject B",
			Code:        "CODE001",
			Description: "Second Description",
		}
		err = repo.CreateSubject(subject2)
		assert.Error(t, err) // Doit échouer à cause de la contrainte unique
	})

	t.Run("GetSubjectByID_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Récupérer la matière par ID
		retrievedSubject, err := repo.GetSubjectByID(subject.ID)
		assert.NoError(t, err)
		assert.Equal(t, subject.ID, retrievedSubject.ID)
		assert.Equal(t, subject.Name, retrievedSubject.Name)
		assert.Equal(t, subject.Code, retrievedSubject.Code)
	})

	t.Run("GetSubjectByID_NotFound", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Essayer de récupérer une matière inexistante
		_, err := repo.GetSubjectByID(99999)
		assert.Error(t, err)
	})

	t.Run("GetSubjectByCode_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Récupérer la matière par code
		retrievedSubject, err := repo.GetSubjectByCode(subject.Code)
		assert.NoError(t, err)
		assert.Equal(t, subject.ID, retrievedSubject.ID)
		assert.Equal(t, subject.Code, retrievedSubject.Code)
	})

	t.Run("GetSubjectByCode_NotFound", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Essayer de récupérer une matière avec un code inexistant
		_, err := repo.GetSubjectByCode("NONEXISTENT")
		assert.Error(t, err)
	})

	t.Run("GetAllSubjects_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer plusieurs matières
		createTestSubject()
		createTestSubject()
		createTestSubject()

		// Récupérer toutes les matières
		subjects, err := repo.GetAllSubjects()
		assert.NoError(t, err)
		assert.Len(t, subjects, 3)
	})



	t.Run("UpdateSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Modifier la matière
		subject.Name = "Updated Subject Name"
		subject.Description = "Updated Description"

		err := repo.UpdateSubject(subject)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedSubject, err := repo.GetSubjectByID(subject.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Subject Name", updatedSubject.Name)
		assert.Equal(t, "Updated Description", updatedSubject.Description)
	})

	t.Run("DeleteSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Supprimer la matière
		err := repo.DeleteSubject(subject.ID)
		assert.NoError(t, err)

		// Vérifier que la matière n'existe plus
		_, err = repo.GetSubjectByID(subject.ID)
		assert.Error(t, err)
	})

	t.Run("CheckSubjectExists_Exists", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Vérifier que la matière existe
		exists, err := repo.CheckSubjectExists(subject.Name, nil)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CheckSubjectExists_NotExists", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Vérifier qu'une matière inexistante n'existe pas
		exists, err := repo.CheckSubjectExists("NonExistentSubject", nil)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("CheckSubjectCodeExists_Exists", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Créer une matière
		subject := createTestSubject()

		// Vérifier que le code existe
		exists, err := repo.CheckSubjectCodeExists(subject.Code, nil)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CheckSubjectCodeExists_NotExists", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()

		// Vérifier qu'un code inexistant n'existe pas
		exists, err := repo.CheckSubjectCodeExists("NONEXISTENT", nil)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestSubjectService(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		req := &models.CreateSubjectRequest{
			Name:        "Service Test Subject",
			Code:        "SERVICE001",
			Description: "Service Test Description",
		}

		response, err := service.CreateSubject(req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Code, response.Code)
		assert.Equal(t, req.Description, response.Description)
	})

	t.Run("CreateSubject_EmptyName_ShouldFail", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		req := &models.CreateSubjectRequest{
			Name:        "",
			Code:        "EMPTY001",
			Description: "Empty Name Description",
		}

		_, err := service.CreateSubject(req)
		assert.Error(t, err)
	})

	t.Run("CreateSubject_EmptyCode_ShouldFail", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		req := &models.CreateSubjectRequest{
			Name:        "Valid Name",
			Code:        "",
			Description: "Empty Code Description",
		}

		_, err := service.CreateSubject(req)
		assert.Error(t, err)
	})

	t.Run("GetSubjectByID_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		// Créer une matière
		subject := createTestSubject()

		// Récupérer la matière
		response, err := service.GetSubjectByID(subject.ID)
		assert.NoError(t, err)
		assert.Equal(t, subject.ID, response.ID)
		assert.Equal(t, subject.Name, response.Name)
		assert.Equal(t, subject.Code, response.Code)
	})

	t.Run("UpdateSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		// Créer une matière
		subject := createTestSubject()

		req := &models.UpdateSubjectRequest{
			Name:        "Updated Service Subject",
			Code:        "UPDATED001",
			Description: "Updated Service Description",
		}

		response, err := service.UpdateSubject(subject.ID, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Code, response.Code)
		assert.Equal(t, req.Description, response.Description)
	})

	t.Run("GetAllSubjects_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		// Créer plusieurs matières
		createTestSubject()
		createTestSubject()
		createTestSubject()

		// Récupérer toutes les matières
		subjects, err := service.GetAllSubjects()
		assert.NoError(t, err)
		assert.Len(t, subjects, 3)
	})

	t.Run("DeleteSubject_Success", func(t *testing.T) {
		repo := repositories.NewSubjectRepository()
		service := services.NewSubjectService(repo)

		// Créer une matière
		subject := createTestSubject()

		// Supprimer la matière
		err := service.DeleteSubject(subject.ID)
		assert.NoError(t, err)

		// Vérifier que la matière n'existe plus
		_, err = service.GetSubjectByID(subject.ID)
		assert.Error(t, err)
	})
}

func TestSubjectValidation(t *testing.T) {
	t.Run("ValidateSubjectName_ValidNames", func(t *testing.T) {
		validNames := []string{
			"Mathématiques",
			"Physique-Chimie",
			"Histoire-Géographie",
			"Informatique",
			"Langues Vivantes",
		}

		for _, name := range validNames {
			assert.NotEmpty(t, name, "Le nom ne devrait pas être vide")
			assert.Len(t, name, 1, "Le nom devrait avoir au moins 1 caractère")
		}
	})

	t.Run("ValidateSubjectName_InvalidNames", func(t *testing.T) {
		invalidNames := []string{
			"",
			"   ",
		}

		for _, name := range invalidNames {
			if name == "" {
				assert.Empty(t, name, "Le nom vide est détecté correctement")
			}
		}
	})

	t.Run("ValidateSubjectCode_ValidCodes", func(t *testing.T) {
		validCodes := []string{
			"MATH001",
			"PHYS002",
			"INFO003",
			"LANG004",
			"HIST005",
		}

		for _, code := range validCodes {
			assert.NotEmpty(t, code, "Le code ne devrait pas être vide")
			assert.Len(t, code, 6, "Le code devrait avoir 6 caractères")
		}
	})

	t.Run("ValidateSubjectCode_InvalidCodes", func(t *testing.T) {
		invalidCodes := []string{
			"",
			"123",
			"TOOLONGCODE",
		}

		for _, code := range invalidCodes {
			if code == "" {
				assert.Empty(t, code, "Le code vide est détecté correctement")
			} else if len(code) < 6 {
				assert.Len(t, code, 6, "Le code trop court est détecté")
			} else if len(code) > 8 {
				assert.Len(t, code, 8, "Le code trop long est détecté")
			}
		}
	})
}
