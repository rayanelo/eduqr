package controllers

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubjectController struct {
	subjectService *services.SubjectService
}

func NewSubjectController(subjectService *services.SubjectService) *SubjectController {
	return &SubjectController{subjectService: subjectService}
}

// GetAllSubjects récupère toutes les matières
func (c *SubjectController) GetAllSubjects(ctx *gin.Context) {
	subjects, err := c.subjectService.GetAllSubjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  subjects,
		"total": len(subjects),
	})
}

// GetSubjectByID récupère une matière par son ID
func (c *SubjectController) GetSubjectByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	subject, err := c.subjectService.GetSubjectByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Matière non trouvée"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subject})
}

// CreateSubject crée une nouvelle matière
func (c *SubjectController) CreateSubject(ctx *gin.Context) {
	var req models.CreateSubjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := c.subjectService.CreateSubject(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": subject})
}

// UpdateSubject met à jour une matière
func (c *SubjectController) UpdateSubject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req models.UpdateSubjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := c.subjectService.UpdateSubject(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subject})
}

// DeleteSubject supprime une matière
func (c *SubjectController) DeleteSubject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	err = c.subjectService.DeleteSubject(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Matière supprimée avec succès"})
}
