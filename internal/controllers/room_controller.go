package controllers

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomService *services.RoomService
}

func NewRoomController(roomService *services.RoomService) *RoomController {
	return &RoomController{roomService: roomService}
}

// GetAllRooms récupère toutes les salles avec filtres
func (c *RoomController) GetAllRooms(ctx *gin.Context) {
	var filter models.RoomFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rooms, err := c.roomService.GetAllRooms(&filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"total": len(rooms),
	})
}

// GetRoomByID récupère une salle par son ID
func (c *RoomController) GetRoomByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	room, err := c.roomService.GetRoomByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Salle non trouvée"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": room})
}

// CreateRoom crée une nouvelle salle
func (c *RoomController) CreateRoom(ctx *gin.Context) {
	var req models.CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if req.IsModular && req.SubRoomsCount < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Une salle modulable doit avoir au moins 2 sous-salles"})
		return
	}

	room, err := c.roomService.CreateRoom(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": room})
}

// UpdateRoom met à jour une salle
func (c *RoomController) UpdateRoom(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var req models.UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := c.roomService.UpdateRoom(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": room})
}

// DeleteRoom supprime une salle
func (c *RoomController) DeleteRoom(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	err = c.roomService.DeleteRoom(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Salle supprimée avec succès"})
}

// GetModularRooms récupère toutes les salles modulables
func (c *RoomController) GetModularRooms(ctx *gin.Context) {
	rooms, err := c.roomService.GetModularRooms()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"total": len(rooms),
	})
}
