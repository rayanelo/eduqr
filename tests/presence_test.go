package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPresenceRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreatePresence_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer les dépendances
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		scannedAt := time.Now()

		presence := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "present",
			ScannedAt: &scannedAt,
		}

		err := repo.CreatePresence(presence)
		assert.NoError(t, err)
		assert.NotZero(t, presence.ID)
		assert.Equal(t, student.ID, presence.StudentID)
		assert.Equal(t, course.ID, presence.CourseID)
		assert.Equal(t, "present", presence.Status)
	})

	t.Run("GetPresenceByID_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer une présence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		presence := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "present",
		}
		testDB.Create(presence)

		// Récupérer la présence par ID
		retrievedPresence, err := repo.GetPresenceByID(presence.ID)
		assert.NoError(t, err)
		assert.Equal(t, presence.ID, retrievedPresence.ID)
		assert.Equal(t, presence.StudentID, retrievedPresence.StudentID)
		assert.Equal(t, presence.CourseID, retrievedPresence.CourseID)
		assert.Equal(t, presence.Status, retrievedPresence.Status)
	})

	t.Run("GetPresenceByID_NotFound", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Essayer de récupérer une présence inexistante
		_, err := repo.GetPresenceByID(99999)
		assert.Error(t, err)
	})

	t.Run("UpdatePresence_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer une présence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		presence := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "present",
		}
		testDB.Create(presence)

		// Modifier la présence
		presence.Status = "late"
		scannedAt := time.Now()
		presence.ScannedAt = &scannedAt

		err := repo.UpdatePresence(presence)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedPresence, err := repo.GetPresenceByID(presence.ID)
		assert.NoError(t, err)
		assert.Equal(t, "late", updatedPresence.Status)
		assert.NotNil(t, updatedPresence.ScannedAt)
	})

	t.Run("DeletePresence_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer une présence
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		presence := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "present",
		}
		testDB.Create(presence)

		// Supprimer la présence
		err := repo.DeletePresence(presence.ID)
		assert.NoError(t, err)

		// Vérifier que la présence n'existe plus
		_, err = repo.GetPresenceByID(presence.ID)
		assert.Error(t, err)
	})

	t.Run("GetPresencesByStudent_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer des présences pour un étudiant
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		presence1 := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "present",
		}
		presence2 := &models.Presence{
			StudentID: student.ID,
			CourseID:  course.ID,
			Status:    "late",
		}
		testDB.Create(presence1)
		testDB.Create(presence2)

		// Récupérer les présences de l'étudiant
		presences, err := repo.GetPresencesByStudent(student.ID)
		assert.NoError(t, err)
		assert.Len(t, presences, 2)

		for _, presence := range presences {
			assert.Equal(t, student.ID, presence.StudentID)
		}
	})

	t.Run("GetPresencesByCourse_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer des présences pour un cours
		student1 := createTestUser("student")
		student2 := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		presence1 := &models.Presence{
			StudentID: student1.ID,
			CourseID:  course.ID,
			Status:    "present",
		}
		presence2 := &models.Presence{
			StudentID: student2.ID,
			CourseID:  course.ID,
			Status:    "absent",
		}
		testDB.Create(presence1)
		testDB.Create(presence2)

		// Récupérer les présences du cours
		presences, err := repo.GetPresencesByCourse(course.ID)
		assert.NoError(t, err)
		assert.Len(t, presences, 2)

		for _, presence := range presences {
			assert.Equal(t, course.ID, presence.CourseID)
		}
	})

	t.Run("GetPresenceStats_Success", func(t *testing.T) {
		repo := repositories.NewPresenceRepository(testDB)

		// Créer des présences avec différents statuts
		student := createTestUser("student")
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Créer des présences avec différents statuts
		statuses := []string{"present", "late", "absent", "present", "late"}
		for _, status := range statuses {
			presence := &models.Presence{
				StudentID: student.ID,
				CourseID:  course.ID,
				Status:    status,
			}
			testDB.Create(presence)
		}

		// Récupérer les statistiques
		stats, err := repo.GetPresenceStats(course.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), stats.PresentStudents)
		assert.Equal(t, int64(2), stats.LateStudents)
		assert.Equal(t, int64(1), stats.AbsentStudents)
	})
}

func TestPresenceValidation(t *testing.T) {
	t.Run("ValidatePresenceStatus_ValidStatuses", func(t *testing.T) {
		validStatuses := []string{
			"present",
			"late",
			"absent",
		}

		for _, status := range validStatuses {
			assert.Contains(t, validStatuses, status, "Le statut doit être valide")
		}
	})

	t.Run("ValidatePresenceStatus_InvalidStatuses", func(t *testing.T) {
		invalidStatuses := []string{
			"",
			"invalid",
			"on_time",
			"missing",
		}

		validStatuses := []string{"present", "late", "absent"}
		for _, status := range invalidStatuses {
			if status != "" {
				assert.NotContains(t, validStatuses, status, "Le statut invalide est détecté")
			}
		}
	})

	t.Run("ValidateScannedAt_ValidTimes", func(t *testing.T) {
		validTimes := []time.Time{
			time.Now(),
			time.Now().Add(-1 * time.Hour),
			time.Now().Add(1 * time.Hour),
		}

		for _, scannedAt := range validTimes {
			assert.NotZero(t, scannedAt, "L'heure de scan ne devrait pas être zéro")
		}
	})

	t.Run("ValidateScannedAt_InvalidTimes", func(t *testing.T) {
		invalidTimes := []time.Time{
			time.Time{}, // Zero time
		}

		for _, scannedAt := range invalidTimes {
			assert.True(t, scannedAt.IsZero(), "L'heure de scan zéro est détectée correctement")
		}
	})
}
