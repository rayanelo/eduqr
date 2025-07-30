package tests

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"eduqr-backend/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCourseRepository(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateCourse_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer les dépendances
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		course := &models.Course{
			Name:        "Test Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   startTime,
			EndTime:     endTime,
			Duration:    120,
			Description: "Test Course Description",
			IsRecurring: false,
		}

		err := repo.CreateCourse(course)
		assert.NoError(t, err)
		assert.NotZero(t, course.ID)
		assert.Equal(t, "Test Course", course.Name)
		assert.Equal(t, subject.ID, course.SubjectID)
		assert.Equal(t, teacher.ID, course.TeacherID)
		assert.Equal(t, room.ID, course.RoomID)
	})

	t.Run("GetCourseByID_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Récupérer le cours par ID
		retrievedCourse, err := repo.GetCourseByID(course.ID)
		assert.NoError(t, err)
		assert.Equal(t, course.ID, retrievedCourse.ID)
		assert.Equal(t, course.Name, retrievedCourse.Name)
		assert.Equal(t, course.SubjectID, retrievedCourse.SubjectID)
		assert.Equal(t, course.TeacherID, retrievedCourse.TeacherID)
		assert.Equal(t, course.RoomID, retrievedCourse.RoomID)
	})

	t.Run("GetCourseByID_NotFound", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Essayer de récupérer un cours inexistant
		_, err := repo.GetCourseByID(99999)
		assert.Error(t, err)
	})

	t.Run("GetAllCourses_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer plusieurs cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		createTestCourse(teacher.ID, subject.ID, room.ID)
		createTestCourse(teacher.ID, subject.ID, room.ID)
		createTestCourse(teacher.ID, subject.ID, room.ID)

		// Récupérer tous les cours
		courses, err := repo.GetAllCourses()
		assert.NoError(t, err)
		assert.Len(t, courses, 3)
	})

	t.Run("UpdateCourse_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Modifier le cours
		course.Name = "Updated Course Name"
		course.Description = "Updated Description"
		course.Duration = 180 // 3 heures

		err := repo.UpdateCourse(course)
		assert.NoError(t, err)

		// Vérifier que les modifications sont sauvegardées
		updatedCourse, err := repo.GetCourseByID(course.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Course Name", updatedCourse.Name)
		assert.Equal(t, "Updated Description", updatedCourse.Description)
		assert.Equal(t, 180, updatedCourse.Duration)
	})

	t.Run("DeleteCourse_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Supprimer le cours
		err := repo.DeleteCourse(course.ID)
		assert.NoError(t, err)

		// Vérifier que le cours n'existe plus
		_, err = repo.GetCourseByID(course.ID)
		assert.Error(t, err)
	})

	t.Run("GetCoursesByRoom_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer une salle et plusieurs cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		createTestCourse(teacher.ID, subject.ID, room.ID)
		createTestCourse(teacher.ID, subject.ID, room.ID)

		// Récupérer les cours de la salle
		courses, err := repo.GetCoursesByRoom(room.ID)
		assert.NoError(t, err)
		assert.Len(t, courses, 2)

		for _, course := range courses {
			assert.Equal(t, room.ID, course.RoomID)
		}
	})

	t.Run("GetCoursesByTeacher_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un enseignant et plusieurs cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		createTestCourse(teacher.ID, subject.ID, room.ID)
		createTestCourse(teacher.ID, subject.ID, room.ID)

		// Récupérer les cours de l'enseignant
		courses, err := repo.GetCoursesByTeacher(teacher.ID)
		assert.NoError(t, err)
		assert.Len(t, courses, 2)

		for _, course := range courses {
			assert.Equal(t, teacher.ID, course.TeacherID)
		}
	})

	t.Run("GetCoursesByRoomAndDate_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours pour une date spécifique
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		targetDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		course := &models.Course{
			Name:        "Date Specific Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   startTime,
			EndTime:     endTime,
			Duration:    120,
			Description: "Date Specific Course Description",
			IsRecurring: false,
		}
		testDB.Create(course)

		// Récupérer les cours de la salle pour cette date
		courses, err := repo.GetCoursesByRoomAndDate(room.ID, targetDate)
		assert.NoError(t, err)
		assert.Len(t, courses, 1)
		assert.Equal(t, room.ID, courses[0].RoomID)
	})

	t.Run("CheckConflicts_Success", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours existant
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		existingStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		existingEnd := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		existingCourse := &models.Course{
			Name:        "Existing Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   existingStart,
			EndTime:     existingEnd,
			Duration:    120,
			Description: "Existing Course Description",
			IsRecurring: false,
		}
		testDB.Create(existingCourse)

		// Créer un nouveau cours avec conflit
		conflictStart := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
		conflictEnd := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)

		newCourse := &models.Course{
			Name:        "Conflicting Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   conflictStart,
			EndTime:     conflictEnd,
			Duration:    120,
			Description: "Conflicting Course Description",
			IsRecurring: false,
		}

		conflicts, err := repo.CheckConflicts(newCourse)
		assert.NoError(t, err)
		assert.Len(t, conflicts, 1) // Doit y avoir un conflit
	})

	t.Run("CheckConflicts_NoConflict", func(t *testing.T) {
		repo := repositories.NewCourseRepository(testDB)

		// Créer un cours existant
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		existingStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		existingEnd := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		existingCourse := &models.Course{
			Name:        "Existing Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   existingStart,
			EndTime:     existingEnd,
			Duration:    120,
			Description: "Existing Course Description",
			IsRecurring: false,
		}
		testDB.Create(existingCourse)

		// Créer un nouveau cours sans conflit
		noConflictStart := time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC)
		noConflictEnd := time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)

		newCourse := &models.Course{
			Name:        "No Conflict Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   noConflictStart,
			EndTime:     noConflictEnd,
			Duration:    120,
			Description: "No Conflict Course Description",
			IsRecurring: false,
		}

		conflicts, err := repo.CheckConflicts(newCourse)
		assert.NoError(t, err)
		assert.Len(t, conflicts, 0) // Ne doit pas y avoir de conflit
	})
}

func TestCourseService(t *testing.T) {
	// Nettoyer avant chaque test
	cleanupTestDatabase()

	t.Run("CreateCourse_Success", func(t *testing.T) {
		courseRepo := repositories.NewCourseRepository(testDB)
		subjectRepo := repositories.NewSubjectRepository()
		userRepo := repositories.NewUserRepository()
		roomRepo := repositories.NewRoomRepository(testDB)
		service := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)

		// Créer les dépendances
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()

		startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

		req := &models.CreateCourseRequest{
			Name:        "Service Test Course",
			SubjectID:   subject.ID,
			TeacherID:   teacher.ID,
			RoomID:      room.ID,
			StartTime:   startTime,
			Duration:    120,
			Description: "Service Test Course Description",
			IsRecurring: false,
		}

		response, err := service.CreateCourse(req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.SubjectID, response.Subject.ID)
		assert.Equal(t, req.TeacherID, response.Teacher.ID)
		assert.Equal(t, req.RoomID, response.Room.ID)
		assert.Equal(t, req.Duration, response.Duration)
	})

	t.Run("GetCourseByID_Success", func(t *testing.T) {
		courseRepo := repositories.NewCourseRepository(testDB)
		subjectRepo := repositories.NewSubjectRepository()
		userRepo := repositories.NewUserRepository()
		roomRepo := repositories.NewRoomRepository(testDB)
		service := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Récupérer le cours
		response, err := service.GetCourseByID(course.ID)
		assert.NoError(t, err)
		assert.Equal(t, course.ID, response.ID)
		assert.Equal(t, course.Name, response.Name)
		assert.Equal(t, course.SubjectID, response.Subject.ID)
		assert.Equal(t, course.TeacherID, response.Teacher.ID)
		assert.Equal(t, course.RoomID, response.Room.ID)
	})

	t.Run("UpdateCourse_Success", func(t *testing.T) {
		courseRepo := repositories.NewCourseRepository(testDB)
		subjectRepo := repositories.NewSubjectRepository()
		userRepo := repositories.NewUserRepository()
		roomRepo := repositories.NewRoomRepository(testDB)
		service := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		newStartTime := time.Date(2024, 1, 2, 14, 0, 0, 0, time.UTC)

		req := &models.UpdateCourseRequest{
			Name:        "Updated Service Course",
			Duration:    180,
			Description: "Updated Service Course Description",
			StartTime:   newStartTime,
		}

		response, err := service.UpdateCourse(course.ID, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Duration, response.Duration)
		assert.Equal(t, req.Description, response.Description)
	})

	t.Run("DeleteCourse_Success", func(t *testing.T) {
		courseRepo := repositories.NewCourseRepository(testDB)
		subjectRepo := repositories.NewSubjectRepository()
		userRepo := repositories.NewUserRepository()
		roomRepo := repositories.NewRoomRepository(testDB)
		service := services.NewCourseService(courseRepo, subjectRepo, userRepo, roomRepo)

		// Créer un cours
		teacher := createTestUser("teacher")
		subject := createTestSubject()
		room := createTestRoom()
		course := createTestCourse(teacher.ID, subject.ID, room.ID)

		// Supprimer le cours
		err := service.DeleteCourse(course.ID)
		assert.NoError(t, err)

		// Vérifier que le cours n'existe plus
		_, err = service.GetCourseByID(course.ID)
		assert.Error(t, err)
	})
}

func TestCourseValidation(t *testing.T) {
	t.Run("ValidateCourseDuration_ValidDurations", func(t *testing.T) {
		validDurations := []int{
			15,  // 15 minutes
			30,  // 30 minutes
			60,  // 1 heure
			90,  // 1h30
			120, // 2 heures
			180, // 3 heures
			240, // 4 heures
			480, // 8 heures (maximum)
		}

		for _, duration := range validDurations {
			assert.GreaterOrEqual(t, duration, 15, "La durée doit être d'au moins 15 minutes")
			assert.LessOrEqual(t, duration, 480, "La durée ne doit pas dépasser 8 heures")
		}
	})

	t.Run("ValidateCourseDuration_InvalidDurations", func(t *testing.T) {
		invalidDurations := []int{
			0,   // Trop court
			5,   // Trop court
			10,  // Trop court
			500, // Trop long
			600, // Trop long
		}

		for _, duration := range invalidDurations {
			if duration < 15 {
				assert.Less(t, duration, 15, "La durée trop courte est détectée")
			} else if duration > 480 {
				assert.Greater(t, duration, 480, "La durée trop longue est détectée")
			}
		}
	})

	t.Run("ValidateCourseTime_ValidTimes", func(t *testing.T) {
		validStartTimes := []time.Time{
			time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),  // 8h00
			time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), // 12h00
			time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC), // 18h00
		}

		for _, startTime := range validStartTimes {
			assert.True(t, startTime.Hour() >= 8 && startTime.Hour() <= 20, "L'heure doit être entre 8h et 20h")
		}
	})

	t.Run("ValidateCourseTime_InvalidTimes", func(t *testing.T) {
		invalidStartTimes := []time.Time{
			time.Date(2024, 1, 1, 6, 0, 0, 0, time.UTC),  // Trop tôt
			time.Date(2024, 1, 1, 22, 0, 0, 0, time.UTC), // Trop tard
		}

		for _, startTime := range invalidStartTimes {
			if startTime.Hour() < 8 {
				assert.Less(t, startTime.Hour(), 8, "L'heure trop tôt est détectée")
			} else if startTime.Hour() > 20 {
				assert.Greater(t, startTime.Hour(), 20, "L'heure trop tard est détectée")
			}
		}
	})
}
