package controllers

import (
	"net/http"
	"strconv"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AbsenceController struct {
	absenceService *services.AbsenceService
}

func NewAbsenceController(absenceService *services.AbsenceService) *AbsenceController {
	return &AbsenceController{
		absenceService: absenceService,
	}
}

// CreateAbsence crée une nouvelle absence
// @Summary Créer une nouvelle absence
// @Description Permet à un étudiant de créer une absence avec justificatif
// @Tags absences
// @Accept json
// @Produce json
// @Param absence body models.CreateAbsenceRequest true "Données de l'absence"
// @Success 201 {object} models.AbsenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /absences [post]
func (c *AbsenceController) CreateAbsence(ctx *gin.Context) {
	var req models.CreateAbsenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'ID de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	absence, err := c.absenceService.CreateAbsence(&req, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, absence)
}

// GetAbsenceByID récupère une absence par son ID
// @Summary Récupérer une absence par ID
// @Description Récupère les détails d'une absence spécifique
// @Tags absences
// @Produce json
// @Param id path int true "ID de l'absence"
// @Success 200 {object} models.AbsenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /absences/{id} [get]
func (c *AbsenceController) GetAbsenceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absence, err := c.absenceService.GetAbsenceByID(uint(id), userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, absence)
}

// GetMyAbsences récupère les absences de l'étudiant connecté
// @Summary Récupérer mes absences
// @Description Récupère les absences de l'étudiant connecté
// @Tags absences
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /absences/my [get]
func (c *AbsenceController) GetMyAbsences(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absences, total, err := c.absenceService.GetAbsencesByStudent(userID.(uint), page, limit, userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  absences,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetTeacherAbsences récupère les absences pour les cours d'un professeur
// @Summary Récupérer les absences des cours du professeur
// @Description Récupère les absences pour les cours du professeur connecté
// @Tags absences
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /absences/teacher [get]
func (c *AbsenceController) GetTeacherAbsences(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absences, total, err := c.absenceService.GetAbsencesByTeacher(userID.(uint), page, limit, userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  absences,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetAllAbsences récupère toutes les absences (admin seulement)
// @Summary Récupérer toutes les absences
// @Description Récupère toutes les absences (réservé aux admins)
// @Tags absences
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/absences [get]
func (c *AbsenceController) GetAllAbsences(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	// Récupérer le rôle de l'utilisateur connecté
	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absences, total, err := c.absenceService.GetAllAbsences(page, limit, userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  absences,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetAbsencesWithFilters récupère les absences avec filtres
// @Summary Récupérer les absences avec filtres
// @Description Récupère les absences avec filtres avancés
// @Tags absences
// @Accept json
// @Produce json
// @Param filters body models.AbsenceFilterRequest true "Filtres à appliquer"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /absences/filter [post]
func (c *AbsenceController) GetAbsencesWithFilters(ctx *gin.Context) {
	var filters models.AbsenceFilterRequest
	if err := ctx.ShouldBindJSON(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absences, total, err := c.absenceService.GetAbsencesWithFilters(&filters, userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  absences,
		"total": total,
		"page":  filters.Page,
		"limit": filters.Limit,
	})
}

// ReviewAbsence valide ou rejette une absence
// @Summary Valider ou rejeter une absence
// @Description Permet à un professeur ou admin de valider/rejeter une absence
// @Tags absences
// @Accept json
// @Produce json
// @Param id path int true "ID de l'absence"
// @Param review body models.ReviewAbsenceRequest true "Décision de validation"
// @Success 200 {object} models.AbsenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /absences/{id}/review [post]
func (c *AbsenceController) ReviewAbsence(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req models.ReviewAbsenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	absence, err := c.absenceService.ReviewAbsence(uint(id), &req, userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, absence)
}

// DeleteAbsence supprime une absence
// @Summary Supprimer une absence
// @Description Supprime une absence (soft delete)
// @Tags absences
// @Produce json
// @Param id path int true "ID de l'absence"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /absences/{id} [delete]
func (c *AbsenceController) DeleteAbsence(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	err = c.absenceService.DeleteAbsence(uint(id), userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Absence supprimée avec succès"})
}

// GetAbsenceStats récupère les statistiques des absences
// @Summary Récupérer les statistiques des absences
// @Description Récupère les statistiques des absences selon le rôle de l'utilisateur
// @Tags absences
// @Produce json
// @Success 200 {object} models.AbsenceStatsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /absences/stats [get]
func (c *AbsenceController) GetAbsenceStats(ctx *gin.Context) {
	// Récupérer les informations de l'utilisateur connecté
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "utilisateur non authentifié"})
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "rôle utilisateur non défini"})
		return
	}

	stats, err := c.absenceService.GetAbsenceStats(userID.(uint), userRole.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
