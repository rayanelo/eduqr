package controllers

import (
	"net/http"
	"strconv"
	"time"

	"eduqr-backend/internal/models"
	"eduqr-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService *services.EventService
}

func NewEventController(eventService *services.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// @Summary Create a new event
// @Description Create a new calendar event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body models.CreateEventRequest true "Event data"
// @Success 201 {object} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /events [post]
func (c *EventController) CreateEvent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := c.eventService.CreateEvent(userID.(uint), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

// @Summary Get event by ID
// @Description Get a specific event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /events/{id} [get]
func (c *EventController) GetEventByID(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := c.eventService.GetEventByID(uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary Get user events
// @Description Get all events for the current user
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.EventResponse
// @Failure 401 {object} map[string]interface{}
// @Router /events [get]
func (c *EventController) GetUserEvents(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	events, err := c.eventService.GetUserEvents(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

// @Summary Update event
// @Description Update an existing event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param event body models.UpdateEventRequest true "Event update data"
// @Success 200 {object} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /events/{id} [put]
func (c *EventController) UpdateEvent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req models.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := c.eventService.UpdateEvent(uint(id), userID.(uint), &req)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary Delete event
// @Description Delete an event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /events/{id} [delete]
func (c *EventController) DeleteEvent(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	err = c.eventService.DeleteEvent(uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Get events by date range
// @Description Get events within a specific date range
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} models.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /events/range [get]
func (c *EventController) GetEventsByDateRange(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	events, err := c.eventService.GetEventsByDateRange(userID.(uint), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
