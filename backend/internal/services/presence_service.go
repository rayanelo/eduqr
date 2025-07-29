package services

import (
	"crypto/rand"
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type PresenceService struct {
	presenceRepo *repositories.PresenceRepository
	courseRepo   *repositories.CourseRepository
	userRepo     *repositories.UserRepository
}

func NewPresenceService(presenceRepo *repositories.PresenceRepository, courseRepo *repositories.CourseRepository, userRepo *repositories.UserRepository) *PresenceService {
	return &PresenceService{
		presenceRepo: presenceRepo,
		courseRepo:   courseRepo,
		userRepo:     userRepo,
	}
}

// GenerateQRCode génère un QR code pour un cours
func (s *PresenceService) GenerateQRCode(courseID uint) (string, error) {
	// Vérifier que le cours existe
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return "", fmt.Errorf("cours non trouvé: %v", err)
	}

	// Vérifier que le cours est en cours ou va bientôt commencer
	now := time.Now()
	if now.Before(course.StartTime.Add(-15 * time.Minute)) {
		return "", fmt.Errorf("le QR code ne peut être généré que 15 minutes avant le début du cours")
	}

	if now.After(course.EndTime) {
		return "", fmt.Errorf("le cours est déjà terminé")
	}

	// Générer un token unique pour le QR code
	token, err := s.generateUniqueToken()
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token: %v", err)
	}

	// Créer les données du QR code
	qrData := map[string]interface{}{
		"course_id": courseID,
		"token":     token,
		"timestamp": now.Unix(),
	}

	// Encoder en JSON puis en base64
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'encodage des données: %v", err)
	}

	qrCodeData := base64.URLEncoding.EncodeToString(jsonData)
	return qrCodeData, nil
}

// ValidateQRCode valide un QR code et retourne les informations du cours
func (s *PresenceService) ValidateQRCode(qrCodeData string) (*models.QRCodeInfo, error) {
	// Décoder le QR code
	jsonData, err := base64.URLEncoding.DecodeString(qrCodeData)
	if err != nil {
		return nil, fmt.Errorf("QR code invalide")
	}

	var qrData map[string]interface{}
	err = json.Unmarshal(jsonData, &qrData)
	if err != nil {
		return nil, fmt.Errorf("QR code invalide")
	}

	// Extraire les données
	courseIDFloat, ok := qrData["course_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("QR code invalide")
	}
	courseID := uint(courseIDFloat)

	// Récupérer le cours
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("cours non trouvé")
	}

	// Vérifier que le cours est en cours
	now := time.Now()
	isValid := now.After(course.StartTime) && now.Before(course.EndTime)

	// Créer les informations du QR code
	qrInfo := &models.QRCodeInfo{
		CourseID:    course.ID,
		CourseName:  course.Name,
		SubjectName: course.Subject.Name,
		TeacherName: fmt.Sprintf("%s %s", course.Teacher.FirstName, course.Teacher.LastName),
		RoomName:    course.Room.Name,
		StartTime:   course.StartTime,
		EndTime:     course.EndTime,
		QRCodeData:  qrCodeData,
		IsValid:     isValid,
	}

	return qrInfo, nil
}

// ScanQRCode traite le scan d'un QR code par un étudiant
func (s *PresenceService) ScanQRCode(qrCodeData string, studentID uint) (*models.Presence, error) {
	// Valider le QR code
	qrInfo, err := s.ValidateQRCode(qrCodeData)
	if err != nil {
		return nil, err
	}

	if !qrInfo.IsValid {
		return nil, fmt.Errorf("le QR code n'est plus valide (cours terminé ou pas encore commencé)")
	}

	// Vérifier que l'utilisateur est un étudiant
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return nil, fmt.Errorf("étudiant non trouvé")
	}

	if student.Role != models.RoleEtudiant {
		return nil, fmt.Errorf("seuls les étudiants peuvent scanner les QR codes")
	}

	// Vérifier si la présence existe déjà
	existingPresence, err := s.presenceRepo.GetPresenceByStudentAndCourse(studentID, qrInfo.CourseID)
	if err == nil {
		// La présence existe déjà, vérifier si elle a déjà été scannée
		if existingPresence.ScannedAt != nil {
			return nil, fmt.Errorf("vous avez déjà scanné ce QR code")
		}
	}

	// Déterminer le statut selon l'heure de scan
	now := time.Now()
	var status string
	if now.Before(qrInfo.StartTime.Add(15 * time.Minute)) {
		status = models.StatusPresent
	} else if now.Before(qrInfo.StartTime.Add(30 * time.Minute)) {
		status = models.StatusLate
	} else {
		status = models.StatusAbsent
	}

	// Créer ou mettre à jour la présence
	var presence *models.Presence
	if existingPresence != nil {
		// Mettre à jour la présence existante
		existingPresence.Status = status
		existingPresence.ScannedAt = &now
		err = s.presenceRepo.UpdatePresence(existingPresence)
		presence = existingPresence
	} else {
		// Créer une nouvelle présence
		presence = &models.Presence{
			StudentID: studentID,
			CourseID:  qrInfo.CourseID,
			Status:    status,
			ScannedAt: &now,
		}
		err = s.presenceRepo.CreatePresence(presence)
	}

	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'enregistrement de la présence: %v", err)
	}

	// Récupérer la présence avec les relations
	return s.presenceRepo.GetPresenceByID(presence.ID)
}

// GetQRCodeInfo récupère les informations d'un QR code pour affichage
func (s *PresenceService) GetQRCodeInfo(courseID uint) (*models.QRCodeInfo, error) {
	// Récupérer le cours
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("cours non trouvé")
	}

	// Vérifier que le cours est en cours ou va bientôt commencer
	now := time.Now()
	isValid := now.After(course.StartTime.Add(-15*time.Minute)) && now.Before(course.EndTime)

	// Générer le QR code si valide
	var qrCodeData string
	if isValid {
		qrCodeData, err = s.GenerateQRCode(courseID)
		if err != nil {
			return nil, err
		}
	}

	// Créer les informations du QR code
	qrInfo := &models.QRCodeInfo{
		CourseID:    course.ID,
		CourseName:  course.Name,
		SubjectName: course.Subject.Name,
		TeacherName: fmt.Sprintf("%s %s", course.Teacher.FirstName, course.Teacher.LastName),
		RoomName:    course.Room.Name,
		StartTime:   course.StartTime,
		EndTime:     course.EndTime,
		QRCodeData:  qrCodeData,
		IsValid:     isValid,
	}

	return qrInfo, nil
}

// GetPresenceStats récupère les statistiques de présence pour un cours
func (s *PresenceService) GetPresenceStats(courseID uint) (*models.PresenceStatsResponse, error) {
	return s.presenceRepo.GetPresenceStats(courseID)
}

// GetPresencesByCourse récupère toutes les présences d'un cours
func (s *PresenceService) GetPresencesByCourse(courseID uint) ([]models.Presence, error) {
	return s.presenceRepo.GetPresencesByCourse(courseID)
}

// GetPresencesByStudent récupère toutes les présences d'un étudiant
func (s *PresenceService) GetPresencesByStudent(studentID uint) ([]models.Presence, error) {
	return s.presenceRepo.GetPresencesByStudent(studentID)
}

// GetPresencesWithFilters récupère les présences avec filtres
func (s *PresenceService) GetPresencesWithFilters(filters map[string]interface{}, page, limit int) ([]models.Presence, int64, error) {
	return s.presenceRepo.GetPresencesWithFilters(filters, page, limit)
}

// CreatePresenceForAllStudents crée des enregistrements de présence pour tous les étudiants d'un cours
func (s *PresenceService) CreatePresenceForAllStudents(courseID uint) error {
	return s.presenceRepo.CreatePresenceForAllStudents(courseID)
}

// generateUniqueToken génère un token unique pour le QR code
func (s *PresenceService) generateUniqueToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CanViewQRCode vérifie si l'utilisateur peut voir le QR code
func (s *PresenceService) CanViewQRCode(userID uint, courseID uint) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	// Les admins et super admins peuvent voir tous les QR codes
	if user.Role == models.RoleAdmin || user.Role == models.RoleSuperAdmin {
		return true, nil
	}

	// Les professeurs peuvent voir les QR codes de leurs cours
	if user.Role == models.RoleProfesseur {
		course, err := s.courseRepo.GetCourseByID(courseID)
		if err != nil {
			return false, err
		}
		return course.TeacherID == userID, nil
	}

	return false, nil
}

// CanRegenerateQRCode vérifie si l'utilisateur peut régénérer un QR code
func (s *PresenceService) CanRegenerateQRCode(userID uint, courseID uint) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	// Les admins et super admins peuvent régénérer tous les QR codes
	if user.Role == models.RoleAdmin || user.Role == models.RoleSuperAdmin {
		return true, nil
	}

	// Les professeurs peuvent régénérer les QR codes de leurs cours
	if user.Role == models.RoleProfesseur {
		course, err := s.courseRepo.GetCourseByID(courseID)
		if err != nil {
			return false, err
		}
		return course.TeacherID == userID, nil
	}

	return false, nil
}
