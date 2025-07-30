package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsenceRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateAbsence_Success", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Créer les dépendances
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		absence := &models.Absence{
			StudentID:     student.ID,
			CourseID:      course.ID,
			Justification: "Maladie",
			Status:        "pending",
		}

		err := repo.CreateAbsence(absence)
		assert.NoError(t, err)
		assert.NotZero(t, absence.ID)
		assert.Equal(t, student.ID, absence.StudentID)
		assert.Equal(t, course.ID, absence.CourseID)
		assert.Equal(t, "pending", absence.Status)
	})

	t.Run("GetAbsenceByID_Success", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Créer une absence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		absence := &models.Absence{
			StudentID:     student.ID,
			CourseID:      course.ID,
			Justification: "Test Absence",
			Status:        "pending",
		}
		testDB.Create(absence)

		// Récupérer l'absence par ID
		retrievedAbsence, err := repo.GetAbsenceByID(absence.ID)
		assert.NoError(t, err)
		assert.Equal(t, absence.ID, retrievedAbsence.ID)
		assert.Equal(t, absence.StudentID, retrievedAbsence.StudentID)
		assert.Equal(t, absence.CourseID, retrievedAbsence.CourseID)
		assert.Equal(t, absence.Status, retrievedAbsence.Status)
	})

	t.Run("GetAbsenceByID_NotFound", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Essayer de récupérer une absence inexistante
		_, err := repo.GetAbsenceByID(99999)
		assert.Error(t, err)
	})

	t.Run("UpdateAbsence_Success", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Créer une absence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		absence := &models.Absence{
			StudentID:     student.ID,
			CourseID:      course.ID,
			Justification: "Test Absence",
			Status:        "pending",
		}
		testDB.Create(absence)

		// Modifier l'absence
		absence.Status = "approved"
		absence.Justification = "Updated Justification"

		err := repo.UpdateAbsence(absence)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedAbsence, err := repo.GetAbsenceByID(absence.ID)
		assert.NoError(t, err)
		assert.Equal(t, "approved", updatedAbsence.Status)
		assert.Equal(t, "Updated Justification", updatedAbsence.Justification)
	})

	t.Run("DeleteAbsence_Success", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Créer une absence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		absence := &models.Absence{
			StudentID:     student.ID,
			CourseID:      course.ID,
			Justification: "Test Absence",
			Status:        "pending",
		}
		testDB.Create(absence)

		// Supprimer l'absence
		err := repo.DeleteAbsence(absence.ID)
		assert.NoError(t, err)

		// Vérifier que l'absence n'existe plus
		_, err = repo.GetAbsenceByID(absence.ID)
		assert.Error(t, err)
	})

	t.Run("GetAbsenceStats_Success", func(t *testing.T) {
		repo := repositories.NewAbsenceRepository(testDB)

		// Créer des absences avec différents statuts
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Créer des absences avec différents statuts
		statuses := []string{"pending", "approved", "rejected", "pending", "approved"}
		for _, status := range statuses {
			absence := &models.Absence{
				StudentID:     student.ID,
				CourseID:      course.ID,
				Justification: "Absence " + status,
				Status:        status,
			}
			testDB.Create(absence)
		}

		// Récupérer les statistiques
		stats, err := repo.GetAbsenceStats()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), stats.TotalAbsences)
		assert.Equal(t, int64(2), stats.PendingAbsences)
		assert.Equal(t, int64(2), stats.ApprovedAbsences)
		assert.Equal(t, int64(1), stats.RejectedAbsences)
	})
}

func TestAbsenceValidation(t *testing.T) {
	t.Run("ValidateAbsenceStatus_ValidStatuses", func(t *testing.T) {
		validStatuses := []string{
			"pending",
			"approved",
			"rejected",
		}

		for _, status := range validStatuses {
			assert.Contains(t, validStatuses, status, "Le statut doit être valide")
		}
	})

	t.Run("ValidateAbsenceStatus_InvalidStatuses", func(t *testing.T) {
		invalidStatuses := []string{
			"",
			"invalid",
			"cancelled",
			"completed",
		}

		validStatuses := []string{"pending", "approved", "rejected"}
		for _, status := range invalidStatuses {
			if status != "" {
				assert.NotContains(t, validStatuses, status, "Le statut invalide est détecté")
			}
		}
	})

	t.Run("ValidateJustification_ValidJustifications", func(t *testing.T) {
		validJustifications := []string{
			"Maladie",
			"Rendez-vous médical",
			"Problème de transport",
			"Événement familial",
			"Autre raison valide",
		}

		for _, justification := range validJustifications {
			assert.NotEmpty(t, justification, "La justification ne devrait pas être vide")
			assert.Len(t, justification, 1, "La justification devrait avoir au moins 1 caractère")
		}
	})

	t.Run("ValidateJustification_InvalidJustifications", func(t *testing.T) {
		invalidJustifications := []string{
			"",
			"   ",
		}

		for _, justification := range invalidJustifications {
			if justification == "" {
				assert.Empty(t, justification, "La justification vide est détectée correctement")
			}
		}
	})
}
