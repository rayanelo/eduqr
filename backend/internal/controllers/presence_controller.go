package controllers

import (
	"net/http"
	"strconv"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type PresenceController struct {
	presenceService *services.PresenceService
}

func NewPresenceController(presenceService *services.PresenceService) *PresenceController {
	return &PresenceController{
		presenceService: presenceService,
	}
}

// ScanQRCode traite le scan d'un QR code par un étudiant
func (pc *PresenceController) ScanQRCode(c *gin.Context) {
	var req models.ScanQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Récupérer l'ID de l'étudiant depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	studentID := userID.(uint)

	// Traiter le scan du QR code
	presence, err := pc.presenceService.ScanQRCode(req.QRCodeData, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Présence enregistrée avec succès",
		"presence": presence.ToPresenceResponse(),
	})
}

// GetQRCodeInfo récupère les informations d'un QR code pour affichage
func (pc *PresenceController) GetQRCodeInfo(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Vérifier les permissions
	canView, err := pc.presenceService.CanViewQRCode(userID.(uint), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des permissions"})
		return
	}

	if !canView {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas les permissions pour voir ce QR code"})
		return
	}

	// Récupérer les informations du QR code
	qrInfo, err := pc.presenceService.GetQRCodeInfo(uint(courseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": qrInfo})
}

// RegenerateQRCode régénère un QR code
func (pc *PresenceController) RegenerateQRCode(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Vérifier les permissions
	canRegenerate, err := pc.presenceService.CanRegenerateQRCode(userID.(uint), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des permissions"})
		return
	}

	if !canRegenerate {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas les permissions pour régénérer ce QR code"})
		return
	}

	// Régénérer le QR code
	qrCodeData, err := pc.presenceService.GenerateQRCode(uint(courseID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "QR code régénéré avec succès",
		"qr_code_data": qrCodeData,
	})
}

// GetPresenceStats récupère les statistiques de présence pour un cours
func (pc *PresenceController) GetPresenceStats(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Vérifier les permissions
	canView, err := pc.presenceService.CanViewQRCode(userID.(uint), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des permissions"})
		return
	}

	if !canView {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas les permissions pour voir ces statistiques"})
		return
	}

	// Récupérer les statistiques
	stats, err := pc.presenceService.GetPresenceStats(uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des statistiques"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetPresencesByCourse récupère toutes les présences d'un cours
func (pc *PresenceController) GetPresencesByCourse(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Vérifier les permissions
	canView, err := pc.presenceService.CanViewQRCode(userID.(uint), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des permissions"})
		return
	}

	if !canView {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas les permissions pour voir ces présences"})
		return
	}

	// Récupérer les présences
	presences, err := pc.presenceService.GetPresencesByCourse(uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des présences"})
		return
	}

	// Convertir en réponses
	var responses []models.PresenceResponse
	for _, presence := range presences {
		responses = append(responses, presence.ToPresenceResponse())
	}

	c.JSON(http.StatusOK, gin.H{"presences": responses})
}

// GetMyPresences récupère les présences de l'étudiant connecté
func (pc *PresenceController) GetMyPresences(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Récupérer les présences
	presences, err := pc.presenceService.GetPresencesByStudent(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des présences"})
		return
	}

	// Convertir en réponses
	var responses []models.PresenceResponse
	for _, presence := range presences {
		responses = append(responses, presence.ToPresenceResponse())
	}

	c.JSON(http.StatusOK, gin.H{"presences": responses})
}

// GetPresencesWithFilters récupère les présences avec filtres (pour les admins)
func (pc *PresenceController) GetPresencesWithFilters(c *gin.Context) {
	// Récupérer les paramètres de pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Construire les filtres
	filters := make(map[string]interface{})

	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := strconv.ParseUint(courseIDStr, 10, 32); err == nil {
			filters["course_id"] = uint(courseID)
		}
	}

	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		if studentID, err := strconv.ParseUint(studentIDStr, 10, 32); err == nil {
			filters["student_id"] = uint(studentID)
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	// Récupérer les présences avec filtres
	presences, total, err := pc.presenceService.GetPresencesWithFilters(filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des présences"})
		return
	}

	// Convertir en réponses
	var responses []models.PresenceResponse
	for _, presence := range presences {
		responses = append(responses, presence.ToPresenceResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreatePresenceForAllStudents crée des enregistrements de présence pour tous les étudiants d'un cours
func (pc *PresenceController) CreatePresenceForAllStudents(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte d'authentification
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Vérifier les permissions
	canView, err := pc.presenceService.CanViewQRCode(userID.(uint), uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des permissions"})
		return
	}

	if !canView {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas les permissions pour cette action"})
		return
	}

	// Créer les enregistrements de présence
	err = pc.presenceService.CreatePresenceForAllStudents(uint(courseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création des enregistrements de présence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enregistrements de présence créés avec succès"})
}
