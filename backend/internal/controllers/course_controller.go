package controllers

import (
	"net/http"
	"strconv"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	courseService *services.CourseService
}

func NewCourseController(courseService *services.CourseService) *CourseController {
	return &CourseController{
		courseService: courseService,
	}
}

// GetAllCourses récupère tous les cours
func (c *CourseController) GetAllCourses(ctx *gin.Context) {
	courses, err := c.courseService.GetAllCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    courses,
	})
}

// GetCourseByID récupère un cours par son ID
func (c *CourseController) GetCourseByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	course, err := c.courseService.GetCourseByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cours non trouvé"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    course,
	})
}

// CreateCourse crée un nouveau cours
func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var req models.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	course, err := c.courseService.CreateCourse(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    course,
		"message": "Cours créé avec succès",
	})
}

// UpdateCourse met à jour un cours existant
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req models.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	course, err := c.courseService.UpdateCourse(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    course,
		"message": "Cours mis à jour avec succès",
	})
}

// DeleteCourse supprime un cours
func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	err = c.courseService.DeleteCourse(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du cours"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cours supprimé avec succès",
	})
}

// GetCoursesByDateRange récupère les cours dans une plage de dates
func (c *CourseController) GetCoursesByDateRange(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Les dates de début et de fin sont requises"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format de date de début invalide (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format de date de fin invalide (YYYY-MM-DD)"})
		return
	}

	courses, err := c.courseService.GetCoursesByDateRange(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    courses,
	})
}

// GetCoursesByRoom récupère les cours d'une salle
func (c *CourseController) GetCoursesByRoom(ctx *gin.Context) {
	roomIDStr := ctx.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de salle invalide"})
		return
	}

	// Récupérer le paramètre de date optionnel
	dateStr := ctx.Query("date")
	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format de date invalide. Utilisez YYYY-MM-DD"})
			return
		}
	} else {
		targetDate = time.Now()
	}

	courses, err := c.courseService.GetCoursesByRoomAndDate(uint(roomID), targetDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    courses,
	})
}

// GetCoursesByTeacher récupère les cours d'un enseignant
func (c *CourseController) GetCoursesByTeacher(ctx *gin.Context) {
	teacherIDStr := ctx.Param("teacherId")
	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID d'enseignant invalide"})
		return
	}

	courses, err := c.courseService.GetCoursesByTeacher(uint(teacherID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des cours"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    courses,
	})
}

// CheckConflicts vérifie les conflits pour un cours
func (c *CourseController) CheckConflicts(ctx *gin.Context) {
	var req models.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	conflicts, err := c.courseService.CheckConflicts(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des conflits"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":       true,
		"data":          conflicts,
		"has_conflicts": len(conflicts) > 0,
	})
}

// CheckConflictsForUpdate vérifie les conflits pour la modification d'un cours
func (c *CourseController) CheckConflictsForUpdate(ctx *gin.Context) {
	courseIDStr := ctx.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cours invalide"})
		return
	}

	var req models.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	conflicts, err := c.courseService.CheckConflictsForUpdate(uint(courseID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des conflits"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":       true,
		"data":          conflicts,
		"has_conflicts": len(conflicts) > 0,
	})
}
