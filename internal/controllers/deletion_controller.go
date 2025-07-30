package controllers

import (
	"net/http"
	"strconv"

	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type DeletionController struct {
	deletionService *services.DeletionService
}

func NewDeletionController(deletionService *services.DeletionService) *DeletionController {
	return &DeletionController{
		deletionService: deletionService,
	}
}

// DeleteUser supprime un utilisateur selon les règles de sécurité
func (c *DeletionController) DeleteUser(ctx *gin.Context) {
	// Récupérer l'ID de l'utilisateur à supprimer
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID d'utilisateur invalide"})
		return
	}

	// Récupérer le rôle de l'utilisateur actuel
	currentUserRole, exists := ctx.Get("current_user_role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "non autorisé"})
		return
	}

	// Récupérer le paramètre de confirmation
	confirmWithCourses := ctx.Query("confirm_with_courses") == "true"

	// Effectuer la suppression
	response, err := c.deletionService.DeleteUser(uint(id), currentUserRole.(string), confirmWithCourses)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Success {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusConflict, response)
	}
}

// DeleteRoom supprime une salle selon les règles de sécurité
func (c *DeletionController) DeleteRoom(ctx *gin.Context) {
	// Récupérer l'ID de la salle à supprimer
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de salle invalide"})
		return
	}

	// Effectuer la suppression
	response, err := c.deletionService.DeleteRoom(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Success {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusConflict, response)
	}
}

// DeleteSubject supprime une matière selon les règles de sécurité
func (c *DeletionController) DeleteSubject(ctx *gin.Context) {
	// Récupérer l'ID de la matière à supprimer
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de matière invalide"})
		return
	}

	// Effectuer la suppression
	response, err := c.deletionService.DeleteSubject(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Success {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusConflict, response)
	}
}

// DeleteCourse supprime un cours selon les règles de sécurité
func (c *DeletionController) DeleteCourse(ctx *gin.Context) {
	// Récupérer l'ID du cours à supprimer
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	// Récupérer le paramètre pour la suppression récurrente
	deleteRecurring := ctx.Query("delete_recurring") == "true"

	// Effectuer la suppression
	response, err := c.deletionService.DeleteCourse(uint(id), deleteRecurring)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.Success {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusConflict, response)
	}
}
