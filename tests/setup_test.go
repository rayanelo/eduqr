package tests

import (
	"eduqr-backend/config"
	"eduqr-backend/internal/database"
	"eduqr-backend/internal/models"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain configure l'environnement de test
func TestMain(m *testing.M) {
	// Configuration pour les tests
	os.Setenv("DB_NAME", "eduqr_test_db")
	os.Setenv("DB_USER", "eduqr_user")
	os.Setenv("DB_PASSWORD", "eduqr_password")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_EXPIRATION", "1h")

	// Charger la configuration
	cfg := config.LoadConfig()

	// Se connecter à la base de données de test
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	testDB = database.GetDB()

	// Nettoyer et migrer la base de données de test
	if err := setupTestDatabase(); err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}

	// Exécuter les tests
	code := m.Run()

	// Nettoyer après les tests
	if err := cleanupTestDatabase(); err != nil {
		log.Printf("Failed to cleanup test database: %v", err)
	}

	// Fermer la connexion
	database.CloseDB()

	os.Exit(code)
}

// setupTestDatabase prépare la base de données de test
func setupTestDatabase() error {
	// Supprimer toutes les tables existantes
	tables := []string{
		"audit_logs",
		"presences",
		"absences",
		"courses",
		"subjects",
		"rooms",
		"users",
	}

	for _, table := range tables {
		if err := testDB.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			return err
		}
	}

	// Auto-migrer les modèles
	models := []interface{}{
		&models.User{},
		&models.Room{},
		&models.Subject{},
		&models.Course{},
		&models.Absence{},
		&models.Presence{},
		&models.AuditLog{},
	}

	for _, model := range models {
		if err := testDB.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestDatabase nettoie la base de données de test
func cleanupTestDatabase() error {
	tables := []string{
		"audit_logs",
		"presences",
		"absences",
		"courses",
		"subjects",
		"rooms",
		"users",
	}

	for _, table := range tables {
		if err := testDB.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			return err
		}
	}

	return nil
}

// createTestUser crée un utilisateur de test
func createTestUser(role string) *models.User {
	user := &models.User{
		Email:     "test-" + role + "@eduqr.com",
		FirstName: "Test",
		LastName:  role,
		Password:  "$2a$10$testpassword",
		Role:      role,
		Phone:     "+1234567890",
		Address:   "Test Address",
	}

	testDB.Create(user)
	return user
}

// createTestRoom crée une salle de test
func createTestRoom() *models.Room {
	room := &models.Room{
		Name:      "Test Room",
		Building:  "Test Building",
		Floor:     "1st Floor",
		IsModular: false,
	}

	testDB.Create(room)
	return room
}

// createTestSubject crée une matière de test
func createTestSubject() *models.Subject {
	subject := &models.Subject{
		Name:        "Test Subject",
		Code:        "TEST001",
		Description: "Test Description",
	}

	testDB.Create(subject)
	return subject
}

// createTestCourse crée un cours de test
func createTestCourse(teacherID, subjectID, roomID uint) *models.Course {
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	course := &models.Course{
		Name:        "Test Course",
		TeacherID:   teacherID,
		SubjectID:   subjectID,
		RoomID:      roomID,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    120,
		Description: "Test Course Description",
		IsRecurring: false,
	}

	testDB.Create(course)
	return course
}
